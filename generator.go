package stubs

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// StubMode for generating data
type StubMode uint64

// Has returns true when this mode has the provided flag configured
func (s StubMode) Has(m StubMode) bool {
	return s&m != 0
}

const (
	// Invalid produces a stub which is invalid for a random validation
	Invalid StubMode = 1 << iota
	// InvalidRequired produces a stub which is invalid for required
	InvalidRequired
	// InvalidMaximum produces a stub which is invalid for maximum
	InvalidMaximum
	// InvalidMinimum produces a stub which is invalid for minimum
	InvalidMinimum
	// InvalidMaxLength produces a stub which is invalid for max length
	InvalidMaxLength
	// InvalidMinLength produces a stub which is invalid for min length
	InvalidMinLength
	// InvalidPattern produces a stub which is invalid for pattern
	InvalidPattern
	// InvalidMaxItems produces a stub which is invalid for max items
	InvalidMaxItems
	// InvalidMinItems produces a stub which is invalid for min items
	InvalidMinItems
	// InvalidUniqueItems produces a stub which is invalid for unique items
	InvalidUniqueItems
	// InvalidMultipleOf produces a stub which is invalid for multiple of
	InvalidMultipleOf
	// InvalidEnum produces a stub which is invalid for enum
	InvalidEnum

	// Valid is the default value and generates valid data
	Valid StubMode = 0
)

// Generator generates a stub for a descriptor.
// A descriptor can either be a parameter, response header or json schema
type Generator struct {
	Language string
}

// Generate a stub into the opts.Target
func (s *Generator) Generate(key string, descriptor interface{}) (interface{}, error) {

	switch desc := descriptor.(type) {
	case *spec.Parameter:
		return s.GenParameter(key, desc)
	case spec.Parameter:
		return s.GenParameter(key, &desc)
	case *spec.Header:
		return s.GenHeader(key, desc)
	case spec.Header:
		return s.GenHeader(key, &desc)
	case *spec.Schema:
		return s.GenSchema(key, desc)
	case spec.Schema:
		return s.GenSchema(key, &desc)
	default:
		return nil, fmt.Errorf("%T is unsupported for Generator", descriptor)
	}
}

// GenParameter generates a random value for a parameter
func (s *Generator) GenParameter(key string, param *spec.Parameter) (interface{}, error) {
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	gopts, err := paramGenOpts(key, param)
	if err != nil {
		return nil, err
	}

	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for parameter [%s]", param.Name)
	}

	return datagen(gopts)
}

// GenHeader generates a random value for a header
func (s *Generator) GenHeader(key string, header *spec.Header) (interface{}, error) {
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	gopts, err := headerGenOpts(key, header)
	if err != nil {
		return nil, err
	}

	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for header [%s]", key)
	}

	return datagen(gopts)
}

// GenSchema generates a random value for a schema
func (s *Generator) GenSchema(key string, schema *spec.Schema) (interface{}, error) {
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	gopts, err := schemaGenOpts(key, true, schema)
	if err != nil {
		return nil, err
	}

	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for header [%s]", key)
	}

	return datagen(gopts)
}
