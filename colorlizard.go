package main

import (
	"net/http"
	"os"

	"github.com/color-lizard/config"

	"github.com/color-lizard/internal/controller"

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

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
	log.Error().Err(err).Msg("Exited")

}
