package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"git.target.com/StoreDataMovement/sdputil/pkg/statsd"
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Endpoint struct {
	QueryPath   string `json:"query_path"`
	QueryString string `json:"query_string"`
	Multiple    bool   `json:"multiple"`
	PrimaryKey	bool   `json:"primary_key"`
	ReturnArray bool   `json:"return_array"`
}

const callHeaderHQ = "X-sdp-result-data-source"
const callValueHQ = "hq_API"
const callValueStore = "store_API"

//AutobahnAPICall Calls Autobahn API
func AutobahnAPICall(url string, reqParamMap url.Values) ([]byte, int, error) {

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err == nil {
		req.URL.RawQuery = reqParamMap.Encode()
	}
	client := &http.Client{}

	// Vars for getting data and status from API call
	var status int
	var data []byte
	log.Info().Msgf("Api re-router calling : %s", url)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Error().Err(err).Msg("Unable to Query autobahn consumer API.Falling back To HQ")
		return nil, http.StatusNotFound, err
	} else {
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			//Read the response
			data, err = ioutil.ReadAll(resp.Body)
			status = resp.StatusCode
			if err != nil || !checkBodyHasContent(data) {
				log.Error().Err(err).Msg("Error reading body.")
				return nil, http.StatusNotFound, errors.New("Error reading body")
			}
		} else {
			log.Info().Msg("No data from autobahn consumer api")
			return nil, http.StatusNotFound, errors.New("No data from autobahn consumer api")
		}
	}

	return data, status, nil
}

func checkBodyHasContent(body []byte) bool {
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return false
	}
	value, ok := jsonParsed.Path("search_meta_data.total_results").Data().(float64)
	_, isTcinPresent := jsonParsed.Path("tcin").Data().(string)
	if(!isTcinPresent){
		_, isTcinPresent = jsonParsed.Path("tcin").Data().(float64)
	}

	if  (ok && value == 0.0) || (!ok && !isTcinPresent) {
		return false
	}
	return true
}

//ReadMappings file which contains mapping of incoming endpoints and autobahn endpoints
func ReadMappings(endpoints *map[string]Endpoint) error {

	mappingFile := "./mapping.json"
	if file := os.Getenv("MAPPING_FILE"); file != "" {
		mappingFile = file
	}
	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	d.UseNumber()
	err = d.Decode(&endpoints)
	if err != nil {
		return err
	}
	return nil
}

//GetHQData calls CallHQAPI function to get results from the HQ
func GetHQData(c *gin.Context, hq HQCaller, cfg *Config) ([]byte, int, error) {

	startTime := time.Now()
	c.Header(callHeaderHQ, callValueHQ)

	data, status, err := hq.CallHQAPI(cfg, c)
	if err != nil {
		//Error from HQ API call
		log.Error().Err(err).Msg("Error from AutoBahn->CallHQAPI()")
		sdpstatsd.WriteErrorCount(errors.New("HQCallFail"))
		return nil, http.StatusNotFound, err
	}

	// got data. record the timing for our HQ call
	duration := time.Since(startTime)
	url := cfg.HQFallbackURL
	query := "TCIN"
	if c.Param("tcin") == "" {
		query = fmt.Sprint(reflect.ValueOf(c.Request.URL.Query()).MapKeys()[0])
	}
	sdpstatsd.WriteAPITime(status, url, "HQ_"+strings.ToUpper(query), c.Request.Method, duration, sdpstatsd.ClientTransaction)

	return data, status, nil
}
