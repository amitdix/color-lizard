package main

import (
	"net/http"

	"config"

	"internal/controller"

	"github.com/rs/zerolog/log"
)

func main() {
	ready := true
	// Read the ENV variables
	var endpointMap map[string]config.Endpoint
	err := config.ReadMockEndpointsData(&endpointMap)
	if err != nil {
		log.Fatal().Msg("Unable To Read Mappings")
	}
	router := controller.GetRouter(endpointMap, &ready)

	http.ListenAndServe(":8881", router)
	log.Error().Err(err).Msg("Exited")

}
