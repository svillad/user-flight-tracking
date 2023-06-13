package api

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/volume/service/user-flight-tracking/controllers"
	"github.com/volume/service/user-flight-tracking/gateways"
	"github.com/volume/service/user-flight-tracking/mediators"
)

// Routes prepares the mux router to be served
func Routes() http.Handler {
	// initialize controllers
	flightTrackerController := generateControllers()

	router := mux.NewRouter()

	// routes
	router.HandleFunc("/calculate", flightTrackerController.GetPath).Methods(http.MethodPost)

	return cors.AllowAll().Handler(router)
}

// generateControllers constructs the needed controller with dependency injected
func generateControllers() controllers.FlightTracker {
	// ------------------------ flightTracker ------------------------
	flightTrackerGateway, _ := gateways.NewFlightTracker(log.WithField("gateway", "FlightTracker"))
	flightTrackerMediator, _ := mediators.NewFlightTracker(log.WithField("mediator", "FlightTracker"), flightTrackerGateway)
	flightTrackerController, _ := controllers.NewFlightTracker(
		log.WithField("controller", "FlightTracker"),
		flightTrackerMediator,
	)

	return flightTrackerController
}
