package validator

// Validator instance
type Validator struct {
	*JSONSchemaValidator
	*StructValidator
}

// NewValidator instance
func NewValidator(jsonSchemaRootPath string) *Validator {
	return &Validator{
		JSONSchemaValidator: NewJSONSchemaValidator(jsonSchemaRootPath),
		StructValidator:     NewStructValidator(),
	}
}
