package translators

import (
	"github.com/volume/service/user-flight-tracking/dto"
	"github.com/volume/service/user-flight-tracking/models"
)

// PathDTOtoModel converts a DTO object into a model object, and returns it.
func PathDTOtoModel(path dto.Path) models.PathResponse {
	var fullPath []string
	for _, flight := range path.Flights {
		fullPath = append(fullPath, flight.Name)
	}

	return models.PathResponse{
		Start: path.Flights[0].Name,
		End:   path.Flights[len(path.Flights)-1].Name,
		Path:  fullPath,
	}
}
