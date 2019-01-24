package sdpstatsd

import (
	"strconv"
	"time"
)

// TransactionType indicates the direction of the API call
type TransactionType uint8

// ServerTransaction is when someone calls us
const ServerTransaction TransactionType = 0

// ClientTransaction is when we call someone else
const ClientTransaction TransactionType = 1

// WriteAPITime should be used for writing API measurements in a predefined way
// - status : return status code
// - url : path portion of the URL being called, without any variable query parameters
// - method : POST, GET, OPTION, etc
// - durationMS : duration of the call. This will be converted to milliseconds; make sure the right units are passed in!
// - transactionType : "server" : endpoing timing provided to our clients : "client" calling some URL e.g. HQ
func WriteAPITime(status int, url string, endpoint string, method string, durationMS time.Duration, transactionType TransactionType) {
	customTags := make(map[string]string)
	duration := durationMS.Nanoseconds() / 1000000
	customTags["status"] = strconv.Itoa(status)
	if transactionType == ServerTransaction {
		customTags["type"] = "server"
	} else {
		customTags["type"] = "client"
	}
	customTags["url"] = url
	customTags["endpoint"] = endpoint
	customTags["method"] = method
	customTags["durationBucket"] = bucketTime(duration)
	sdw.writeTime("API", duration, customTags)
}
