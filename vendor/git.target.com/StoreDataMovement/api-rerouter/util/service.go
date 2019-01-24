package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CallAutobahnApi(key string, cfg *Config, reqParamMap url.Values, endpoints map[string]Endpoint, context *gin.Context, hq HQCaller) ([]byte, int, error) {
	value := reqParamMap.Get(key)
	autobahnQuery := make(url.Values)
	var url string
	if endpoints[key].PrimaryKey == false {
		if _, ok := endpoints[key]; !ok {
			log.Error().Msg("No Mapping Found")
			return nil, http.StatusNotFound, errors.New("No Mapping Found")
		}
		autobahnQuery.Add(endpoints[key].QueryString, value)
		if endpoints[key].Multiple {
			queryList := strings.Split(value, ",")
			if len(queryList) > 1 {
				return GetHQData(context, hq, cfg)
			}
		}
		url = cfg.AutobahnBasePath + endpoints[key].QueryPath
	}else {
		_, isKeyInPath := context.Params.Get(key)
		if !isKeyInPath{
			url = cfg.AutobahnBasePath + endpoints[key].QueryPath + context.Query(key)
		}else{
			url = cfg.AutobahnBasePath + endpoints[key].QueryPath + context.Params.ByName(key)
		}
	}

	body, statuscode, err := AutobahnAPICall(url, autobahnQuery)
	if err != nil {
		return GetHQData(context, hq, cfg)
	}

	var data map[string]interface{}
	error := json.Unmarshal(body, &data)


	if error == nil {
		if val, ok := data["records"]; ok {
			body, error = json.Marshal(val.([]interface{}))
		}
	}

	if endpoints[key].ReturnArray {
		//encode with [] if the output format is an array
		//attribute ReturnArray should be false for any other endpoints
		body = []byte("["+string(body)+"]")
	}

	return body, statuscode, err
}

