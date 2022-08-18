// Package server implements server for collecting metrics from agent cleints.
package server

import (
	"context"
	"flag"
	"net"

	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/andrei-cloud/go-devops/internal/config"
	"github.com/andrei-cloud/go-devops/internal/encrypt"
	"github.com/andrei-cloud/go-devops/internal/interceptors"
	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/andrei-cloud/go-devops/internal/router"
	"github.com/andrei-cloud/go-devops/internal/storage/filestore"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/andrei-cloud/go-devops/internal/storage/persistent"

	pb "github.com/andrei-cloud/go-devops/internal/proto"
)

var (
	cfg        config.ServerConfig
	configPath = flag.String("config", "", "path to config file")
)

type server struct {
	r      *chi.Mux
	s      *http.Server
	g      *grpc.Server
	gl     net.Listener
	repo   repo.Repository
	f      filestore.Filestore
	key    []byte
	subnet *net.IPNet
}

func init() {
	flag.StringVar(configPath, "c", "", "path to config file")

	addressPtr := flag.String("a", "localhost:8080", "server address format: host:port")
	restorePtr := flag.Bool("r", true, "restore previous values")
	intervalPtr := flag.Duration("i", 30*time.Second, "interval to store metrics")
	filePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file path to store metrics")
	keyPtr := flag.String("k", "", "secret key")
	dsnPtr := flag.String("d", "", "database connection string")
	debugPtr := flag.Bool("debug", false, "sets log level to debug")
	cryptokeyPtr := flag.String("cyptokey", "", "path to private key file")
	subnetPtr := flag.String("t", "", "trusted subnet in CIDR format")
	grpcPtr := flag.Bool("grpc", false, "enable grpc communication")

	flag.Parse()
	cfg = config.ServerConfig{}
	if configPath != nil && *configPath != "" {
		config.ReadConfigFile(*configPath, cfg)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal().AnErr("Parse", err).Msg("init")
	}
	if cfg.Address == "" {
		cfg.Address = *addressPtr
	}
	if cfg.Interval == 0 {
		cfg.Interval = *intervalPtr
	}
	if cfg.FilePath == "" {
		cfg.FilePath = *filePtr
	}
	if !cfg.Restore {
		cfg.Restore = false
	} else {
		cfg.Restore = cfg.Restore || *restorePtr
	}
	if cfg.Key == "" {
		cfg.Key = *keyPtr
	}
	if cfg.Dsn == "" {
		cfg.Dsn = *dsnPtr
	}
	if cfg.CryptoKey == "" {
		cfg.CryptoKey = *cryptokeyPtr
	}
	if cfg.Subnet == "" {
		cfg.Subnet = *subnetPtr
	}

	if !cfg.Grpc {
		cfg.Grpc = *grpcPtr
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugPtr {
		cfg.Debug = true
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("DEBUG LEVEL IS ENABLED")
	}
}

// NewServer - sreates new server instance with all ingected dependencies.
func NewServer() *server {
	var (
		decr encrypt.Decrypter
		err  error
	)

	srv := server{}
	srv.repo = inmem.New()

	if cfg.Key != "" {
		srv.key = []byte(cfg.Key)
	}

	if cfg.Dsn != "" {
		log.Debug().Msg("Database is used as Storage")
		srv.repo = persistent.NewDB(cfg.Dsn)
		if srv.repo == nil {
			log.Fatal().Msg("Failed to connect to DB")
		}
	} else if cfg.FilePath != "" {
		log.Debug().Msg("Faile is used as Storage")
		srv.f = filestore.NewFileStorage(cfg.FilePath)
	}

	if cfg.CryptoKey != "" {
		decr = encrypt.New(cfg.CryptoKey)
	}

	_, srv.subnet, err = net.ParseCIDR(cfg.Subnet)
	if err != nil {
		log.Error().AnErr("ParseCIDR", err).Msg("NewServer")
		srv.subnet = nil
	}

	srv.r = router.SetupRouter(srv.repo, srv.key, decr)

	if cfg.Debug {
		srv.r = router.WithPPROF(srv.r)
	}

	srv.s = &http.Server{
		Addr:           cfg.Address,
		Handler:        srv.r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.Grpc {
		if cfg.CryptoKey != "" {
			encoding.RegisterCodec(encrypt.Encodec{Dec: decr})
		}
		srv.gl, err = net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatal().AnErr("Listen", err).Msg("Failed to listen port :9090")
		}

		srv.g = grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.CheckIP(srv.subnet)))
		pb.RegisterMetricsServer(srv.g, NewMetricsServer(srv))
	}

	return &srv
}

// Run - non blocking function starting up the server.
func (srv *server) Run(ctx context.Context) {
	if cfg.Dsn == "" && cfg.FilePath != "" {
		if cfg.Restore {
			if err := srv.f.Restore(srv.repo); err != nil {
				log.Error().AnErr("Restore", err).Msg("Run")
			}
		}

		storeTicker := time.NewTicker(cfg.Interval)

		go func(ctx context.Context) {
			for {
				select {
				case <-storeTicker.C:
					if err := srv.f.Store(srv.repo); err != nil {
						log.Error().AnErr("Store", err).Msg("Run")
					}
				case <-ctx.Done():
					storeTicker.Stop()
					return
				}
			}
		}(ctx)
	}

	log.Info().Msgf("HTTP server listening on: %v", cfg.Address)
	go srv.s.ListenAndServe()

	if cfg.Grpc && srv.gl != nil {
		log.Info().Msgf("gRPC server listening on: :9090")
		go srv.g.Serve(srv.gl)
	}
}

// Shutdown - blocking function waiting signal to shutdown the server
// signals to shutdown server:
//
//	os.Interrupt
//	syscall.SIGTERM
//	syscall.SIGQUIT
//
// Server will be forcefuly stopped after shutdown Timeout.
func (srv *server) Shutdown(ctx context.Context) {
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Shutdown)
	defer cancel()
	if err := srv.s.Shutdown(ctx); err != nil {
		log.Error().AnErr("Shutdown", err).Msg("Shutdown")
	}

	if srv.g != nil {
		srv.g.GracefulStop()
		if err := srv.gl.Close(); err != nil {
			log.Error().AnErr("listener close", err).Msg("Shutdown")
		}
	}

	if srv.f != nil && cfg.FilePath != "" {
		if err := srv.f.Store(srv.repo); err != nil {
			log.Error().AnErr("Store", err).Msg("Shutdown")
		}
	}

	if srv.repo != nil {
		if err := srv.repo.Close(); err != nil {
			log.Error().AnErr("Close", err).Msg("Shutdown")
		}
	}
}
