package mediators

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/gateways"
	"github.com/volume/service/user-flight-tracking/models"
)

// FlightTracker specifies the methods to get flights
type FlightTracker interface {
	GetFlightsPath(ctx context.Context, req models.PathRequest) (dto.Path, error)
}

// flightTracker is the concrete implementation of the FlightTracker interface
type flightTracker struct {
	Logger               *log.Entry
	FlightTrackerGateway gateways.FlightTracker
}

// NewFlightTracker returns a new instance of FlightTracker mediator
func NewFlightTracker(log *log.Entry, flightTrackerGateway gateways.FlightTracker) (FlightTracker, error) {
	switch {
	case log == nil:
		return nil, errors.New("logger")
	case flightTrackerGateway == nil:
		return nil, errors.New("flightTrackerGateway")
	}

	return &flightTracker{
		Logger:               log,
		FlightTrackerGateway: flightTrackerGateway,
	}, nil
}

// GetFlightsPath returns a flights path
func (m *flightTracker) GetFlightsPath(ctx context.Context, req models.PathRequest) (dto.Path, error) {
	var path dto.Path

	path, err := m.FlightTrackerGateway.GetFlightsPath(ctx, req)
	if err != nil {
		return dto.Path{}, err
	}

	return path, nil
}
