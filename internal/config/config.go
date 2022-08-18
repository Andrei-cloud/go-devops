package config

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Config - type for agent configuration.
type AgentConfig struct {
	Address   string        `json:"address" env:"ADDRESS"`                 // address of metric server
	Key       string        `env:"KEY"`                                    // key for metrics hashing
	CryptoKey string        `json:"crypto_key" env:"CRYPTO_KEY"`           // key for encryption of metrics
	ReportInt time.Duration `json:"report_interval" env:"REPORT_INTERVAL"` // interval for metrics reporting
	PollInt   time.Duration `json:"poll_interval" env:"POLL_INTERVAL"`     // interval for metrics polling
	IsBulk    bool          // flag to send metrics in bulk
	Debug     bool          // debug flag
	Grpc      bool          `env:"ENABLE_GRPC"` // enable grpc communication
}

// Config - type for server configuration.
type ServerConfig struct {
	Address   string        `env:"ADDRESS"`                          // address server to bind on
	Key       string        `env:"KEY"`                              // key used for hash verifications
	Dsn       string        `env:"DATABASE_DSN"`                     // dadabase connection string
	FilePath  string        `env:"STORE_FILE"`                       // path to the file to store metrics
	CryptoKey string        `env:"CRYPTO_KEY"`                       // key for encryption of metrics
	Shutdown  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"` // time to wait for server shutdown
	Interval  time.Duration `env:"STORE_INTERVAL"`                   // interval store metrics in persistemnt repository
	Restore   bool          `env:"RESTORE" envDefault:"true"`        // restore metrics from file upon server start
	Debug     bool          // debug mode enables additional logging and profile enpoints
	Subnet    string        `env:"TRUSTED_SUBNET"` // trusted subnet for agent
	Grpc      bool          `env:"ENABLE_GRPC"`    // enable grpc communication
}

func ReadConfigFile(path string, c interface{}) {
	var cfg interface{}

	switch x := c.(type) {
	case *AgentConfig:
		cfg = x
	case *ServerConfig:
		cfg = x
	default:
		log.Fatal().Msg("invalid config type")
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().AnErr("ReadFile", err).Msg("reading config file failed")
	}

	s := string(fileBytes)

	re := regexp.MustCompile(`(?im)\/\/[^"\[\]^{^}]+$`)
	s = re.ReplaceAllString(s, "")
	err = json.NewDecoder(strings.NewReader(s)).Decode(cfg)
	if err != nil {
		log.Fatal().AnErr("Decode", err).Msg("decofing json file failed")
	}
}
