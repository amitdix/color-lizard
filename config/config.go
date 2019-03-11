package config

import (
	"encoding/json"
	"os"
)

type Endpoint struct {
	Method   string            `json:"method"`
	Status   int               `json:status`
	Response string            `json:response`
	Headers  map[string]string `json:"headers"`
}

func ReadMockEndpointsData(endpointMap *map[string]Endpoint) error {
	mappingFile := "./config.json"
	if file := os.Getenv("MAPPING_FILE"); file != "" {
		mappingFile = file
	}
	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	d.UseNumber()
	err = d.Decode(&endpointMap)
	if err != nil {
		return err
	}
	return nil
}
