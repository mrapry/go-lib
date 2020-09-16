package validator

import "github.com/mrapry/go-lib/config"

// Validator instance
type Validator struct {
	*JSONSchemaValidator
	*StructValidator
}

// NewValidator instance
func NewValidator() *Validator {
	return &Validator{
		JSONSchemaValidator: NewJSONSchemaValidator(config.BaseEnv().JSONSchemaDir),
		StructValidator:     NewStructValidator(),
	}
}
