package options

import (
	"errors"

	"github.com/buzzfeed/sso/internal/pkg/sessions"
	"github.com/buzzfeed/sso/internal/proxy/providers"
)

var (
	// These error message should be formatted in such a way that is appropriate
	// for display to the end user.
	ErrInvalidEmailAddress = errors.New("Invalid Email Address In Session State")
	ErrValidationError     = errors.New("Error during validation")
)

type Validator interface {
	Validate(*sessions.SessionState) error
}

// RunValidators runs each passed in validator and returns a slice of errors. If an
// empty slice is returned, it can be assumed all passed in validators were successful.
func RunValidators(validators []Validator, session *sessions.SessionState) []error {
	validatorErrors := make([]error, 0, len(validators))

	for _, validator := range validators {
		err := validator.Validate(session)
		if err != nil {
			validatorErrors = append(validatorErrors, err)
		}
	}
	return validatorErrors
}

func RunValidatorsWithGracePeriod(validators []Validator, session *sessions.SessionState) []error {
	validatorErrors := make([]error, 0, len(validators))
	for _, err := range RunValidators(validators, session) {
		if err, ok := err.(*providers.GroupValidationError); ok {
			if err.IsWithinGracePeriod() {
				continue
			}
		}
		validatorErrors = append(validatorErrors, err)
	}
	return validatorErrors
}
