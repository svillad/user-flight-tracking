package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

// PathRequest model
type PathRequest struct {
	Flights [][]string
}

func (pr PathRequest) Validate() error {
	return validation.ValidateStruct(&pr,
		validation.Field(&pr.Flights,
			validation.Required,
			validation.Each(validation.Length(2, 2).Error("each flght must contain exactly 2 airports")),
			validation.Each(validation.Each(validation.Required, validation.Length(3, 3).Error("each airport must have exactly 3 characters"))),
			validation.Each(validation.Each(validation.Match(regexp.MustCompile(`^\S+$`)).Error("each airport must not contain spaces"))),
		),
	)
}
