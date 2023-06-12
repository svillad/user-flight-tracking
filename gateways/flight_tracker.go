package gateways

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/models"
)

// FlightTracker specifies the methods to get flights information
type FlightTracker interface {
	GetFlightsPath(ctx context.Context, req models.PathRequest) (dto.Path, error)
}

// flightTracker is the concrete implementation of the FlightTracker interface
type flightTracker struct {
	Logger *log.Entry
}

// NewFlightTracker returns a new instance of FlightTracker gateway
func NewFlightTracker(log *log.Entry) (FlightTracker, error) {
	switch {
	case log == nil:
		return nil, errors.New("logger")
	}

	return &flightTracker{
		Logger: log,
	}, nil
}

// Get returns the path of flights from a specific user
func (m *flightTracker) GetFlightsPath(ctx context.Context, req models.PathRequest) (dto.Path, error) {
	var path dto.Path

	graph := buildGraph(req.Flights)
	startFlight, _, err := findStartAndEndFlights(graph)
	if err != nil {
		return dto.Path{}, err
	}
	path.Flights = findPath(startFlight, path.Flights)
	log.Info(buildStringPath(path))

	return path, nil
}

func buildGraph(pairs [][]string) map[string]*dto.Flight {
	graph := make(map[string]*dto.Flight)

	for _, pair := range pairs {
		source := pair[0]
		destination := pair[1]

		// Create new flight if it does not exist in the graph
		if _, ok := graph[source]; !ok {
			graph[source] = &dto.Flight{Name: source}
		}
		if _, ok := graph[destination]; !ok {
			graph[destination] = &dto.Flight{Name: destination}
		}

		// Create the connections between the flights
		sourceFlight := graph[source]
		destinationFlight := graph[destination]
		sourceFlight.Outgoing = append(sourceFlight.Outgoing, destinationFlight)
		destinationFlight.Incoming = append(destinationFlight.Incoming, sourceFlight)
	}

	return graph
}

func findStartAndEndFlights(graph map[string]*dto.Flight) (*dto.Flight, *dto.Flight, error) {
	var startFlight, endFlight *dto.Flight

	for _, node := range graph {
		if len(node.Incoming) == 0 {
			node.IsStart = true
			startFlight = node
		}
		if len(node.Outgoing) == 0 {
			node.IsEnd = true
			endFlight = node
		}
	}

	if startFlight == nil {
		return nil, nil, fmt.Errorf("no initial flight found")
	}

	if endFlight == nil {
		return nil, nil, fmt.Errorf("no final flight found")
	}

	checkInFlights(startFlight)

	// Check circular flights
	if startFlight == endFlight {
		return nil, nil, fmt.Errorf("a circular flight was found between flights: %s", startFlight.Name)
	}

	// Check disconnections
	disconnectedFlights := make([]string, 0)
	for _, node := range graph {
		if !node.Visited {
			disconnectedFlights = append(disconnectedFlights, node.Name)
		}
	}

	if len(disconnectedFlights) > 0 {
		return nil, nil, fmt.Errorf("disconnections detected between flights: %v", disconnectedFlights)
	}

	return startFlight, endFlight, nil
}

func checkInFlights(node *dto.Flight) {
	node.Visited = true

	for _, neighbor := range node.Outgoing {
		if !neighbor.Visited {
			checkInFlights(neighbor)
		}
	}
}

func buildStringPath(path dto.Path) string {
	var p []string
	for _, flight := range path.Flights {
		p = append(p, flight.Name)
	}

	// Join the elements of pathStrings using commas
	return fmt.Sprintf("[%s]=>[%s,%s]", strings.Join(p, ","), p[0], p[len(p)-1])
}

func findPath(node *dto.Flight, path []*dto.Flight) []*dto.Flight {
	path = append(path, node)

	if node.IsEnd {
		return path
	}

	for _, nextFlight := range node.Outgoing {
		nextPath := findPath(nextFlight, append([]*dto.Flight(nil), path...))
		if nextPath != nil {
			return nextPath
		}
	}

	return nil
}
