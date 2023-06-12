package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/volume/service/user-flight-tracking/controllers/translators"
	"github.com/volume/service/user-flight-tracking/mediators"
	"github.com/volume/service/user-flight-tracking/models"
)

// Service defines the methods for flight
type FlightTracker interface {
	GetPath(w http.ResponseWriter, r *http.Request)
}

// flightTracker defines the components for the controller
type flightTracker struct {
	Logger                *log.Entry
	FlightTrackerMediator mediators.FlightTracker
}

// NewFlightTracker returns a new instance of FlightTracker controller
func NewFlightTracker(log *log.Entry, flightTrackerMediator mediators.FlightTracker) (FlightTracker, error) {
	switch {
	case log == nil:
		return nil, errors.New("logger")
	case flightTrackerMediator == nil:
		return nil, errors.New("flightTrackerMediator")
	}

	return &flightTracker{
		Logger:                log,
		FlightTrackerMediator: flightTrackerMediator,
	}, nil
}

// GetPath retrieves flight path from the backend
func (c *flightTracker) GetPath(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	// Decodes the JSON data from the request body into an instance of the `PathRequest` structure
	var request models.PathRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		c.Logger.WithError(err).Error("error decoding JSON")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		c.Logger.WithError(err).Error("error validating request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	path, err := c.FlightTrackerMediator.GetFlightsPath(context.Background(), request)
	if err != nil {
		c.Logger.WithError(err).Error("internal server error")
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(translators.PathDTOtoModel(path)); err != nil {
		c.Logger.WithError(err).Error("error encoding JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
