package sdpstatsd

// WriteErrorCount is for recordoing error occurrences
func WriteErrorCount(err error) {
	sdw.writeCount("Error", 1, map[string]string{"kind": "", "func": "", "msg": err.Error()})
}
