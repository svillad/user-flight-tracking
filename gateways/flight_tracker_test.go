package gateways_test

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
	"github.com/volume/service/user-flight-tracking/models"
)

func TestGateways_NewFlightTracker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger = log.NewEntry(nil)
	)

	type args struct {
		logger *log.Entry
	}
	tests := []struct {
		name      string
		args      args
		wantError error
	}{
		{
			name: "should_return_success",
			args: args{
				logger: logger,
			},
			wantError: nil,
		},
		{
			name: "should_return_error_when_the_logger_is_nil",
			args: args{
				logger: nil,
			},
			wantError: errors.New("logger"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gateways.NewFlightTracker(tt.args.logger)
			if err != nil {
				assert.Equal(t, tt.wantError.Error(), err.Error())
			}
		})
	}
}

func TestGateways_GetFlightsPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger = log.NewEntry(log.New())
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

		req := models.PathRequest{
			Flights: [][]string{
				{"IND", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
		}

		g, err := gateways.NewFlightTracker(logger)
		require.NoError(t, err)

		resp, err := g.GetFlightsPath(context.Background(), req)

		assert.Equal(t, len(wantedPath.Flights), len(resp.Flights))
		assert.Equal(t, wantedPath.Flights[0].Name, resp.Flights[0].Name)
		assert.Equal(t, wantedPath.Flights[4].Name, resp.Flights[4].Name)
		assert.NilError(t, err)
	})

	t.Run("should_return_path", func(t *testing.T) {
		wantedPath := dto.Path{}

		req := models.PathRequest{
			Flights: [][]string{
				{"IND", "SFO"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
		}

		g, err := gateways.NewFlightTracker(logger)
		require.NoError(t, err)

		resp, err := g.GetFlightsPath(context.Background(), req)

		assert.Equal(t, len(wantedPath.Flights), len(resp.Flights))
		assert.Error(t, err, "no initial flight found")
	})

	t.Run("should_return_path", func(t *testing.T) {
		wantedPath := dto.Path{}

		req := models.PathRequest{
			Flights: [][]string{
				{"XXX", "EWR"},
				{"SFO", "ATL"},
				{"GSO", "IND"},
				{"ATL", "GSO"},
			},
		}

		g, err := gateways.NewFlightTracker(logger)
		require.NoError(t, err)

		resp, err := g.GetFlightsPath(context.Background(), req)

		assert.Equal(t, len(wantedPath.Flights), len(resp.Flights))
		assert.Error(t, err, "disconnections detected between flights: [XXX EWR]")
	})
}
