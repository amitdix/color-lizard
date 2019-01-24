package sdpstatsd

import (
	"strconv"

	bjson "git.target.com/StoreDataMovement/sdputil/pkg/json"
	"github.com/rs/zerolog/log"
)

// WriteKafkaMetricBlob parses kafak library metrics and writes them in statsD form
func WriteKafkaMetricBlob(statsJSON string, consumerID string) {
	parseLags(statsJSON, consumerID)
}

// Kafka Metrics Utility Stuff

type kafkaMetric struct {
	Topic, Group string
	Values       map[string]interface{}
}

func makeKafkaMetric(topic string, group string, partition int, mType string, value interface{}) kafkaMetric {
	k := kafkaMetric{topic, group, make(map[string]interface{})}
	k.Values["partition"] = partition
	k.Values["type"] = mType
	k.Values["value"] = value
	return k
}

func parseLags(ks string, consumerID string) {
	failMsg := "Kafka Stats: couldn't parse stats json"
	stats, err := bjson.Parse(ks)
	if err != nil {
		log.Error().Msg(failMsg)
		return
	}
	topics, err := stats.GetJson("topics")
	if err != nil {
		log.Error().Msg(failMsg)
		return
	}
	for top := range topics {
		topic, err := topics.GetJson(top)
		if err != nil {
			log.Debug().Msg("couldn't get stats for topic " + top)
			continue
		}
		partitions, err := topic.GetJson("partitions")
		if err != nil {
			log.Debug().Msg("couldn't get partitions for topic " + top)
			continue
		}
		// Also the app offset (last offset passed to app)
		lagTotal, aOffsetTotal, cOffsetTotal := int64(0), int64(0), int64(0)
		for part := range partitions {
			partition, err := partitions.GetJson(part)
			partitionNum, err1 := strconv.Atoi(part)
			fetchState, err2 := partition.GetString("fetch_state")
			if err != nil || err1 != nil || err2 != nil {
				log.Debug().Msg("couldn't get stats for topic " + top)
				continue
			}
			// make sure this is THIS consumer's partition, and we have accurate stats for it.
			if fetchState != "active" {
				continue
			}
			if partitionNum < 0 {
				continue // skip partition "-1"
			}
			lag, err := partition.GetInt64("consumer_lag")
			if err != nil {
				log.Debug().Msg("couldn't get consumer_lag for topic " + top)
			} else {
				lagTotal += lag
			}
			aOffset, err := partition.GetInt64("app_offset")
			if err != nil {
				log.Debug().Msg("couldn't get stats for topic " + top)
			} else {
				aOffsetTotal += aOffset
			}
			cOffset, err := partition.GetInt64("committed_offset")
			if err != nil {
				log.Debug().Msg("couldn't get stats for topic " + top)
			} else {
				cOffsetTotal += cOffset
			}
		}

		// Write out the combined row with all the pieces
		sdw.writeInfluxLine("KafkaConsumer",
			map[string]string{
				"type":           "lag",
				"consumer_group": consumerID,
				"topic":          top,
				"lagBucket":      bucketLag(lagTotal),
			},
			map[string]interface{}{
				"lag":        lagTotal,
				"appOffset":  aOffsetTotal,
				"commOffset": cOffsetTotal,
			})
	}
}
