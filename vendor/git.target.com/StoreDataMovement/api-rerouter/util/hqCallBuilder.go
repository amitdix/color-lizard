package util

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HQCallerLive defines the real call to HQ ItemLoc API
type HQCallerLive struct {
	client *http.Client
	pool   *x509.CertPool
}

// SetupHQCaller gets Target certs to use for outgoing calls to HTTPS endpoints
func SetupHQCaller(cfg *Config) (*HQCallerLive, error) {
	resp, err := http.Get(cfg.CertsURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed getting Target certs bundle")
		return nil, err
	}
	defer resp.Body.Close()
	pemCerts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed getting Target certs bundle")
		return nil, err
	}

	h := HQCallerLive{}
	h.pool = x509.NewCertPool()
	h.pool.AppendCertsFromPEM(pemCerts)
	h.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: h.pool}}}
	log.Info().Msg("Read certs")
	return &h, nil
}

// CallHQAPI retrieves HQ data for a given query parameter
func (h HQCallerLive) CallHQAPI(cfg *Config, c *gin.Context) ([]byte, int, error) {
	if len(cfg.StoreName) < 5 {
		log.Error().Msg("CallHQAPI: Could not find the store number")
		return nil, http.StatusNotFound, errors.New("CallHQAPI: Could not find the store number")
	}
	runes := []rune(cfg.StoreName)
	strStore := string(runes[1:5])

	store, err := strconv.Atoi(strStore)
	if err != nil {
		//Error: Cant find store number
		log.Error().Msg("CallHQAPI: Could not find the store number in alpha to numeric conversion")
		return nil, http.StatusNotFound, errors.New("CallHQAPI: Could not find the store number in alpha to numeric conversion")
	}

	if cfg.HQFallbackURL == "" {
		log.Error().Msg("CallHQAPI: HQ URL Not Set")
		return nil, http.StatusNotFound, errors.New("CallHQAPI: HQ URL Not Set")
	}

	var url string

	//check for path parameter

	if len(c.Param("tcin")) > 0 {
		tcin, err := strconv.Atoi(c.Param("tcin"))
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse tcin from path.")
			return nil, http.StatusBadRequest, errors.New("CallHQAPI: Unable to parse tcin from path")
		}
		url = fmt.Sprintf("%s/%d?location_id=%d&key=%s", cfg.HQFallbackURL, tcin, store, cfg.APIKey)
		log.Info().Msgf("HQAPI call : %s?tcin=%d&location_id=%d", cfg.HQFallbackURL, tcin, store)
	} else {
		url = fmt.Sprint(cfg.HQFallbackURL)
		log.Info().Msgf("HQAPI call : %s", cfg.HQFallbackURL)
	}

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new http request.")
		return nil, http.StatusNotFound, errors.New("CallHQAPI: Unable to create new http request")
	}
	if err == nil && c.Param("tcin") == "" {
		reqParamMap := c.Request.URL.Query()
		reqParamMap["key"] = []string{cfg.APIKey}
		reqParamMap["location_id"] = []string{strconv.Itoa(store)}
		req.URL.RawQuery = reqParamMap.Encode()
	}
	resp, err := h.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Query HQ API.")
		return nil, resp.StatusCode, errors.New("CallHQAPI: Unable to Query HQ API")
	}
	defer resp.Body.Close()

	//Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading body.")
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Info().Msgf("HQ API returned status code %d", resp.StatusCode)
	}
	return body, resp.StatusCode, nil
}
