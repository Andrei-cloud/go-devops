// Package server implements server for collecting metrics from agent cleints.
package server

import (
	"context"
	"flag"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/andrei-cloud/go-devops/internal/encrypt"
	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/andrei-cloud/go-devops/internal/router"
	"github.com/andrei-cloud/go-devops/internal/storage/filestore"
	"github.com/andrei-cloud/go-devops/internal/storage/inmem"
	"github.com/andrei-cloud/go-devops/internal/storage/persistent"
)

var cfg Config

// Config - type for server configuration.
type Config struct {
	Address   string        `env:"ADDRESS"`                          // address server to bind on
	Key       string        `env:"KEY"`                              // key used for hash verifications
	Dsn       string        `env:"DATABASE_DSN"`                     // dadabase connection string
	FilePath  string        `env:"STORE_FILE"`                       // path to the file to store metrics
	CryptoKey string        `env:"CRYPTO_KEY"`                       // key for encryption of metrics
	Shutdown  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"` // time to wait for server shutdown
	Interval  time.Duration `env:"STORE_INTERVAL"`                   // interval store metrics in persistemnt repository
	Restore   bool          `env:"RESTORE" envDefault:"true"`        // restore metrics from file upon server start
	Debug     bool          // debug mode enables additional logging and profile enpoints
}

type server struct {
	r    *chi.Mux
	s    *http.Server
	repo repo.Repository
	f    filestore.Filestore
	key  []byte
}

func init() {
	addressPtr := flag.String("a", "localhost:8080", "server address format: host:port")
	restorePtr := flag.Bool("r", true, "restore previous values")
	intervalPtr := flag.Duration("i", 30*time.Second, "interval to store metrics")
	filePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file path to store metrics")
	keyPtr := flag.String("k", "", "secret key")
	dsnPtr := flag.String("d", "", "database connection string")
	debugPtr := flag.Bool("debug", false, "sets log level to debug")
	cryptokeyPtr := flag.String("cyptokey", "", "path to private key file")

	flag.Parse()
	cfg = Config{}
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

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debugPtr {
		cfg.Debug = true
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("DEBUG LEVEL IS ENABLED")
	}
}

// NewServer - sreates new server instance with all ingected dependencies.
func NewServer() *server {
	var encr encrypt.Encrypter
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
		encr = encrypt.New(cfg.CryptoKey)
	}
	srv.r = router.SetupRouter(srv.repo, srv.key, encr)

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

	log.Info().Msgf("server listening on: %v", cfg.Address)
	go srv.s.ListenAndServe()

}

// Shutdown - blocking function waiting signal to shutdown the server
// signals to shutdown server:
//   os.Interrupt
//   syscall.SIGTERM
//   syscall.SIGQUIT
// Server will be forcefuly stopped after shutdown Timeout.
func (srv *server) Shutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Shutdown)
	defer cancel()
	if err := srv.s.Shutdown(ctx); err != nil {
		log.Error().AnErr("Shutdown", err).Msg("Shutdown")
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
