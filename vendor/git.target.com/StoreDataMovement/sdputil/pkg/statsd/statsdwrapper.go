package sdpstatsd

import (
	"errors"
	"fmt"
	"sync"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	statsd "github.com/smira/go-statsd"
)

// Static singleton version
var sdw statsdwriter
var initOnce sync.Once

// Internal representation of the statsD lib
type clientStatsDInterface interface {
	Timing(stat string, delta int64, tags ...statsd.Tag)
	Incr(stat string, count int64, tags ...statsd.Tag)
	Gauge(stat string, value int64, tags ...statsd.Tag)
	PrecisionTiming(stat string, delta time.Duration, tags ...statsd.Tag)
	Close() error
}

// Error messages
const errMsgInitAlready = "Metrics already initialized!"

// Init initializes static Statsd
// Externally provided tags:
// - location (telegraf)
// - environment (telegraf)
func Init(host string, statsDPort uint, influxPort uint, appTag string, streamTag string, otherDefaultTags map[string]string) (err error) {
	if sdw.initialized {
		return errors.New(errMsgInitAlready)
	}
	initOnce.Do(func() {
		if len(host) == 0 || statsDPort == 0 || influxPort == 0 {
			err = errors.New("Cannot initialize statsd or influx/telegraf with empty hostname and/or port")
		} else {
			statsDAddr := fmt.Sprintf("%s:%d", host, statsDPort)
			influxAddr := fmt.Sprintf("%s:%d", host, influxPort)
			statsDclient := statsd.NewClient(statsDAddr, statsd.MaxPacketSize(1400))
			influxClient, err := influx.NewUDPClient(influx.UDPConfig{Addr: influxAddr})
			if err != nil {
				panic(err.Error())
			}
			otherDefaultTags["app"] = appTag
			otherDefaultTags["stream"] = streamTag
			err = InitWithClient(statsDclient, influxClient, otherDefaultTags)
		}
	})

	return err
}

// InitWithClient initializes with a given client interface (mostly for testing)
func InitWithClient(client clientStatsDInterface, influxClient influx.Client, defaultTags map[string]string) (err error) {
	if sdw.initialized {
		return errors.New(errMsgInitAlready)
	}
	sdw.tags = make(map[string]string)
	for k, v := range defaultTags {
		sdw.tags[k] = CleanTag(v)
	}
	sdw.client = client
	sdw.influxClient = influxClient
	sdw.initialized = true
	return nil
}

// Close resets the metrics wrapper to an uninitialized state
func Close() {
	sdw.close()
}
