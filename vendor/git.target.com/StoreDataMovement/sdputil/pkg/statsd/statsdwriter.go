package sdpstatsd

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	influx "github.com/influxdata/influxdb/client/v2"
)

// Statsdwriter ...
type statsdwriter struct {
	tags         map[string]string
	client       clientStatsDInterface
	influxClient influx.Client
	initialized  bool
}

const msgErrNotInitialized = "Attempting to log metrics without initializing."

// Internal calls, NOT TO BE USED directly!
func (sdw statsdwriter) close() error {
	if !sdw.initialized {
		log.Debug().Msg("Attempting to close without init")
		return nil
	}

	e1 := sdw.client.Close()
	e2 := sdw.influxClient.Close()
	if e1 == nil && e2 == nil {
		return nil
	}

	// one or both close calls failed :(
	var errStrings []string
	errStrings = append(errStrings, e1.Error())
	errStrings = append(errStrings, e2.Error())
	return errors.New(strings.Join(errStrings, ";"))
}

// count
func (sdw statsdwriter) writeCount(key string, value int64, customTags map[string]string) {
	if !sdw.initialized {
		log.Debug().Msg(msgErrNotInitialized)
		return
	}
	tags := addCustomTags(customTags)
	sdw.client.Incr(key, value, tags...)
}

// gauge
func (sdw statsdwriter) writeGauge(key string, value int64, customTags map[string]string) {
	if !sdw.initialized {
		log.Debug().Msg(msgErrNotInitialized)
		return
	}
	tags := addCustomTags(customTags)
	sdw.client.Gauge(key, value, tags...)
}

// time
func (sdw statsdwriter) writeTime(key string, millis int64, customTags map[string]string) {
	if !sdw.initialized {
		log.Debug().Msg(msgErrNotInitialized)
		return
	}
	tags := addCustomTags(customTags)
	sdw.client.Timing(key, millis, tags...)
}

// time with duration?! Not basic statsD format
func (sdw statsdwriter) writeTimewithDuration(key string, millis time.Duration, rate int64, customTags map[string]string) {
	if !sdw.initialized {
		log.Debug().Msg(msgErrNotInitialized)
		return
	}
	tags := addCustomTags(customTags)
	sdw.client.PrecisionTiming(key, millis, tags...)
}

func (sdw statsdwriter) writeInfluxLine(key string, customTags map[string]string, fields map[string]interface{}) {
	if !sdw.initialized {
		log.Debug().Msg(msgErrNotInitialized)
		return
	}

	tagRecs := addCustomTagsString(customTags)
	tags := make(map[string]string)
	for k, v := range tagRecs {
		tags[k] = v
	}
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Precision: "ns",
	})
	if err != nil {
		log.Warn().Msg("Can't make batchpoints collection")
		return
	}

	pt, err := influx.NewPoint(
		key,
		tags,
		fields,
		time.Now())
	if err != nil {
		log.Warn().Msg("Can't make new point")
		return
	}

	bp.AddPoint(pt)

	if err := sdw.influxClient.Write(bp); err != nil {
		log.Warn().Err(err).Msg("Can't write point")
	}
}
