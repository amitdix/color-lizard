package sdpstatsd

import (
	"time"
)

// OperationStatus captures results of a processing step
type OperationStatus uint8

// Possible statuses of operations
const (
	OpOK      OperationStatus = 0
	OpFailure OperationStatus = 1
	OpRetry   OperationStatus = 2
	OpUnknown OperationStatus = 3
)

func (m OperationStatus) String() string {
	names := []string{
		"success",
		"failure",
		"retry",
		"unknown",
	}
	return names[m]
}

// WriteOperation records multiple measurements to capture count, timing and rate for a processing function
func WriteOperation(operationName string, status OperationStatus, duration time.Duration, recCount uint32) {
	tags := map[string]string{"name": operationName, "result": status.String(), "timeBucket": bucketTime(int64(duration / time.Millisecond))}

	// We use 'timings' to get the local statsD aggregation and statistics of our measurements
	sdw.writeTime("OperationCount", int64(recCount), tags)
	sdw.writeTime("OperationTime", int64(duration/time.Millisecond), tags)
	if duration == 0 {
		duration = 1
	}
	sdw.writeTime("OperationRate", int64((time.Duration(recCount)*time.Second)/duration), tags)
}
