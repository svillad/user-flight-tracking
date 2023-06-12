package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/volume/service/user-flight-tracking/controllers"
	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/mediators"
	mock_flightTracker_mediator "github.com/volume/service/user-flight-tracking/mocks/mockmediators"
	"github.com/volume/service/user-flight-tracking/models"
)

func TestController_NewFlightTracker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger       = log.NewEntry(nil)
		mockMediator = mock_flightTracker_mediator.NewMockFlightTracker(ctrl)
	)

	type args struct {
		logger       *log.Entry
		mockMediator mediators.FlightTracker
	}
	tests := []struct {
		name      string
		args      args
		wantError error
	}{
		{
			name: "should_return_success",
			args: args{
				logger:       logger,
				mockMediator: mockMediator,
			},
			wantError: nil,
		},
		{
			name: "should_return_error_when_the_logger_is_nil",
			args: args{
				logger:       nil,
				mockMediator: mockMediator,
			},
			wantError: errors.New("logger"),
		},
		{
			name: "should_return_error_when_the_mediator_is_nil",
			args: args{
				logger:       logger,
				mockMediator: nil,
			},
			wantError: errors.New("flightTrackerMediator"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := controllers.NewFlightTracker(tt.args.logger, tt.args.mockMediator)
			if err != nil {
				assert.Equal(t, tt.wantError.Error(), err.Error())
			}
		})
	}
}

func TestController_GetPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		logger       = log.NewEntry(log.New())
		mockMediator = mock_flightTracker_mediator.NewMockFlightTracker(ctrl)
	)

	t.Run("should_return_path", func(t *testing.T) {
		wantedPath := models.PathResponse{
			Start: "SFO",
			End:   "EWR",
			Path:  []string{"SFO", "ATL", "GSO", "IND", "EWR"},
		}

		path := dto.Path{
			Flights: []*dto.Flight{
				{Name: "SFO"},
				{Name: "ATL"},
				{Name: "GSO"},
				{Name: "IND"},
				{Name: "EWR"},
			},
		}

		mockMediator.EXPECT().GetFlightsPath(gomock.Any(), gomock.Any()).Return(path, nil)

		c, err := controllers.NewFlightTracker(logger, mockMediator)
		require.NoError(t, err)

		jsonBody := `{
			"flights": [
				["IND", "EWR"],
				["SFO", "ATL"],
				["GSO", "IND"],
				["ATL", "GSO"]
			]
		}`

		// Crea un lector a partir de la cadena de texto JSON
		bodyReader := bytes.NewReader([]byte(jsonBody))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/calculate", bodyReader)

		c.GetPath(recorder, request)

		resp := recorder.Result()
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")

		responseBody := models.PathResponse{}
		err = json.Unmarshal(body, &responseBody)
		require.NoError(t, err, "should unmarshal the response wrapper without error")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, wantedPath.Start, responseBody.Start)
		assert.Equal(t, wantedPath.End, responseBody.End)
		assert.Equal(t, len(wantedPath.Path), len(responseBody.Path))
	})

	t.Run("failure_response_when_bad_request", func(t *testing.T) {
		c, err := controllers.NewFlightTracker(logger, mockMediator)
		require.NoError(t, err)

		jsonBody := `{
			"flights": [
				["IND", "EWR",
				["SFO", "ATL"],
				["GSO", "IND"],
				["ATL", "GSO"]
			]
		}`

		// Crea un lector a partir de la cadena de texto JSON
		bodyReader := bytes.NewReader([]byte(jsonBody))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/calculate", bodyReader)

		c.GetPath(recorder, request)

		resp := recorder.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure_response_when_bad_request", func(t *testing.T) {
		c, err := controllers.NewFlightTracker(logger, mockMediator)
		require.NoError(t, err)

		jsonBody := `{
			"flights": [
				["IND"],
				["SFO", "ATL"],
				["GSO", "IND"],
				["ATL", "GSO"]
			]
		}`

		// Crea un lector a partir de la cadena de texto JSON
		bodyReader := bytes.NewReader([]byte(jsonBody))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/calculate", bodyReader)

		c.GetPath(recorder, request)

		resp := recorder.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure_response_when_mediator_retrun_error", func(t *testing.T) {
		c, err := controllers.NewFlightTracker(logger, mockMediator)
		require.NoError(t, err)

		jsonBody := `{
			"flights": [
				["IND", "EWR"],
				["SFO", "ATL"],
				["GSO", "IND"],
				["ATL", "GSO"]
			]
		}`

		mockMediator.EXPECT().GetFlightsPath(gomock.Any(), gomock.Any()).Return(dto.Path{}, errors.New("internal server error"))

		// Crea un lector a partir de la cadena de texto JSON
		bodyReader := bytes.NewReader([]byte(jsonBody))
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/calculate", bodyReader)

		c.GetPath(recorder, request)

		resp := recorder.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
