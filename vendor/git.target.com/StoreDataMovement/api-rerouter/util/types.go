package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Config captures DB and HQ API settings for the app
type Config struct {
	StoreName        string `env:"STORE_NAME" envDefault:"T0003"`
	Port             string `env:"PORT" envDefault:"9000"`
	StatsDHost       string `env:"STATSD_HOST" envDefault:"localhost"`
	StatsDPort       uint   `env:"STATSD_PORT" envDefault:"8125"`
	InfluxPort       uint   `env:"INFLUX_PORT" envDefault:"8094"`
	GinMode          string `env:"GIN_MODE" envDefault:"release"`
	LogLevel         uint   `env:"LOG_LEVEL" envDefault:"1"`
	Version          string `env:"VERSION" envDefault:"0.0.0"`
	APIBasePath      string `env:"API_BASE_PATH" envDefault:"/store_core_items/v1"`
	AutobahnBasePath string `env:"AB_BASE_PATH" envDefault:"http://localhost:5050"`
	AppTag           string `env:"APP_TAG" envDefault:"SDPApp"`
	CertsURL         string `env:"CERTS_URL" envDefault:"http://browserconfig.target.com/tgt-certs/tgt-ca-bundle.crt"`
	HQFallbackURL    string `env:"HQ_Fallback_URL" envDefault:"https://api-internal.target.com/store_items/v1/store_items"`
	APIKey           string `env:"API_KEY" envDefault:"AnAPIKeyGoesHere"`
}

// HQCaller provides item loc data from an HQ api
type HQCaller interface {
	CallHQAPI(cfg *Config, c *gin.Context) ([]byte, int, error)
}

// ErrorMessage Error json
type ErrorMessage struct {
	Error string
	code  string
}

// CreateError Returns error object
func CreateError(message string, code int) ErrorMessage {
	var errorObject ErrorMessage
	errorObject.Error = message
	errorObject.code = strconv.Itoa(code)
	return errorObject
}

// Error provides a way to declace a constant error type
type Error string

func (e Error) Error() string { return string(e) }
