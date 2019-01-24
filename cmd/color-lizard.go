package main

import (
	"git.target.com/StoreDataMovement/color-lizard/config"
	"git.target.com/StoreDataMovement/color-lizard/internal/controller"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main(){
	ready := true
	// Read the ENV variables
	var endpoints []config.Endpoint
	err := config.ReadMockEndpointsData(&endpoints)
	if err != nil {
		log.Fatal().Msg("Unable To Read Mappings")
	}
	router := controller.GetRouter(endpoints, &ready)

	http.ListenAndServe(":8881", router)
	log.Error().Err(err).Msg("Exited")

}

