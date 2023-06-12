package translators_test

import (
	"testing"

	"github.com/volume/service/user-flight-tracking/controllers/translators"
	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/models"
	"gotest.tools/assert"
)

func TestTranslator_PathDTOtoModel(t *testing.T) {
	cases := []struct {
		name      string
		pathDTO   dto.Path
		pathModel models.PathResponse
	}{
		{
			name: "Successful translation",
			pathDTO: dto.Path{
				Flights: []*dto.Flight{
					{Name: "SFO"},
					{Name: "ATL"},
					{Name: "GSO"},
					{Name: "IND"},
					{Name: "EWR"},
				},
			},
			pathModel: models.PathResponse{
				Start: "SFO",
				End:   "EWR",
				Path:  []string{"SFO", "ATL", "GSO", "IND", "EWR"},
			},
		},
	}

	for _, c := range cases {
		response := translators.PathDTOtoModel(c.pathDTO)
		assert.Equal(t, c.pathModel.Start, response.Start)
		assert.Equal(t, c.pathModel.End, response.End)
		assert.Equal(t, c.pathModel.Path[0], response.Path[0])
		assert.Equal(t, c.pathModel.Path[4], response.Path[4])
		assert.Equal(t, len(c.pathModel.Path), len(response.Path))
	}
}
