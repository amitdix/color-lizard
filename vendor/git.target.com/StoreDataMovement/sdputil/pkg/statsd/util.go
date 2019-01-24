package sdpstatsd

import (
	"regexp"

	statsd "github.com/smira/go-statsd"
)

// groups times into buckets for Grafana / Influx to manage all tag values (as strings) in a roughly logarithmic scale
func bucketTime(dur int64) string {
	switch {
	case dur < 50:
		return "a"
	case dur < 100:
		return "b"
	case dur < 500:
		return "c"
	case dur < 1000:
		return "d"
	case dur < 2000:
		return "e"
	case dur < 5000:
		return "f"
	case dur < 10000:
		return "g"
	default:
		return "z"
	}
}

// groups times into buckets for Grafana / Influx to manage all tag values (as strings) in a roughly logarithmic scale
func bucketLag(lag int64) string {
	switch {
	case lag < 100:
		return "a"
	case lag < 1000:
		return "b"
	case lag < 10000:
		return "c"
	case lag < 100000:
		return "d"
	case lag < 500000:
		return "e"
	case lag < 1000000:
		return "f"
	case lag < 2000000:
		return "g"
	case lag < 5000000:
		return "h"
	default:
		return "z"
	}
}

// helper to add the set of custom tags passed in to the list of tags on the measurement, along with the default tags
func addCustomTags(customTags map[string]string) (tags []statsd.Tag) {
	var statsdCustomTags []statsd.Tag
	tagDict := addCustomTagsString(customTags)
	for k, v := range tagDict {
		var tempTag = statsd.StringTag(k, CleanTag(v))
		statsdCustomTags = append(statsdCustomTags, tempTag)
	}
	return statsdCustomTags
}

func addCustomTagsString(customTags map[string]string) map[string]string {
	allTags := make(map[string]string)
	for k, v := range sdw.tags {
		allTags[k] = v
	}
	if customTags != nil {
		for k, v := range customTags {
			allTags[k] = v
		}
	}
	return allTags
}

func CleanTag(s string) string {
	// Strip any query params, replace everything else non-alphanumeric with '_'
	re := regexp.MustCompile(`\?.*`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`[^-A-Za-z0-9/._]`)
	s = re.ReplaceAllString(s, "_")
	return s
}
