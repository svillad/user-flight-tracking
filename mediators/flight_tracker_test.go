package mediators_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/gateways"
	"github.com/volume/service/user-flight-tracking/mediators"
	mock_flightTracker_gateway "github.com/volume/service/user-flight-tracking/mocks/mockgateways"
	"github.com/volume/service/user-flight-tracking/models"
)

func TestMediators_NewFlightTracker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger      = log.NewEntry(nil)
		mockGateway = mock_flightTracker_gateway.NewMockFlightTracker(ctrl)
	)

	type args struct {
		logger      *log.Entry
		mockGateway gateways.FlightTracker
	}
	tests := []struct {
		name      string
		args      args
		wantError error
	}{
		{
			name: "should_return_success",
			args: args{
				logger:      logger,
				mockGateway: mockGateway,
			},
			wantError: nil,
		},
		{
			name: "should_return_error_when_the_logger_is_nil",
			args: args{
				logger:      nil,
				mockGateway: mockGateway,
			},
			wantError: errors.New("logger"),
		},
		{
			name: "should_return_error_when_the_mediator_is_nil",
			args: args{
				logger:      logger,
				mockGateway: nil,
			},
			wantError: errors.New("flightTrackerGateway"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mediators.NewFlightTracker(tt.args.logger, tt.args.mockGateway)
			if err != nil {
				assert.Equal(t, tt.wantError.Error(), err.Error())
			}
		})
	}
}

func TestMediators_GetFlightsPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger      = log.NewEntry(log.New())
		mockGateway = mock_flightTracker_gateway.NewMockFlightTracker(ctrl)
	)

	t.Run("should_return_path", func(t *testing.T) {
		wantedPath := dto.Path{
			Flights: []*dto.Flight{
				{Name: "SFO"},
				{Name: "ATL"},
				{Name: "GSO"},
				{Name: "IND"},
				{Name: "EWR"},
			},
		}

		mockGateway.EXPECT().GetFlightsPath(gomock.Any(), gomock.Any()).Return(wantedPath, nil)

		m, err := mediators.NewFlightTracker(logger, mockGateway)
		require.NoError(t, err)

		resp, err := m.GetFlightsPath(context.Background(), models.PathRequest{})

		assert.Equal(t, len(wantedPath.Flights), len(resp.Flights))
		assert.NilError(t, err)
	})

	t.Run("failure_response_when_gateway_retrun_error", func(t *testing.T) {
		mockGateway.EXPECT().GetFlightsPath(gomock.Any(), gomock.Any()).Return(dto.Path{}, errors.New("internal server error"))

		m, err := mediators.NewFlightTracker(logger, mockGateway)
		require.NoError(t, err)

		resp, err := m.GetFlightsPath(context.Background(), models.PathRequest{})

		assert.Equal(t, 0, len(resp.Flights))
		assert.Error(t, err, "internal server error")
	})
}
