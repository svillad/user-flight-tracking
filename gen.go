package gen

//go:generate mockgen -package mock_flightTracker_controller -destination mocks/mockcontrollers/flightTracker_mock.go github.com/volume/service/user-flight-tracking/controllers FlightTracker
//go:generate mockgen -package mock_flightTracker_mediator -destination mocks/mockmediators/flightTracker_mock.go github.com/volume/service/user-flight-tracking/mediators FlightTracker
//go:generate mockgen -package mock_flightTracker_gateway -destination mocks/mockgateways/flightTracker_mock.go github.com/volume/service/user-flight-tracking/gateways FlightTracker
