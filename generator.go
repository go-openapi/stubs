package stubs

import (
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/mitchellh/mapstructure"
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
	// InvalidMaximum produces a stub which is invalid for required
	InvalidMaximum
	// InvalidMinimum produces a stub which is invalid for required
	InvalidMinimum
	// InvalidMaxLength produces a stub which is invalid for required
	InvalidMaxLength
	// InvalidMinLength produces a stub which is invalid for required
	InvalidMinLength
	// InvalidPattern produces a stub which is invalid for required
	InvalidPattern
	// InvalidMaxItems produces a stub which is invalid for required
	InvalidMaxItems
	// InvalidMinItems produces a stub which is invalid for required
	InvalidMinItems
	// InvalidUniqueItems produces a stub which is invalid for required
	InvalidUniqueItems
	// InvalidMultipleOf produces a stub which is invalid for required
	InvalidMultipleOf
	// InvalidEnum produces a stub which is invalid for required
	InvalidEnum

	// Valid is the default value and generates valid data
	Valid StubMode = 0
)

type simpleOpts struct {
	spec.CommonValidations
	spec.SimpleSchema

	name      string
	args      []interface{}
	fieldName string
	required  bool
	mode      StubMode
}

func (g *simpleOpts) Mode() StubMode {
	return g.mode
}

func (g *simpleOpts) Args() []interface{} {
	return g.args
}

func (g *simpleOpts) CollectionFormat() string {
	return g.SimpleSchema.CollectionFormat
}

func (g *simpleOpts) Name() string {
	return g.name
}
func (g *simpleOpts) FieldName() string {
	return g.fieldName
}
func (g *simpleOpts) Maximum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Maximum), g.CommonValidations.ExclusiveMaximum, g.CommonValidations.Maximum != nil
}
func (g *simpleOpts) Minimum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Minimum), g.CommonValidations.ExclusiveMinimum, g.CommonValidations.Minimum != nil
}
func (g *simpleOpts) MaxLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MaxLength), g.CommonValidations.MaxLength != nil
}
func (g *simpleOpts) MinLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MinLength), g.CommonValidations.MinLength != nil
}
func (g *simpleOpts) Pattern() (string, bool) {
	return g.CommonValidations.Pattern, g.CommonValidations.Pattern != ""
}
func (g *simpleOpts) MaxItems() (int64, bool) {
	mx := g.CommonValidations.MaxItems
	return swag.Int64Value(mx), mx != nil
}
func (g *simpleOpts) MinItems() (int64, bool) {
	mn := g.CommonValidations.MinItems
	return swag.Int64Value(mn), mn != nil
}
func (g *simpleOpts) UniqueItems() bool {
	return g.CommonValidations.UniqueItems
}
func (g *simpleOpts) MultipleOf() (float64, bool) {
	mo := g.CommonValidations.MultipleOf
	return swag.Float64Value(mo), mo != nil
}
func (g *simpleOpts) Enum() ([]interface{}, bool) {
	enm := g.CommonValidations.Enum
	return enm, len(enm) > 0
}
func (g *simpleOpts) Type() string {
	return g.SimpleSchema.Type
}
func (g *simpleOpts) Format() string {
	return g.SimpleSchema.Format
}
func (g *simpleOpts) Items() (GeneratorOpts, error) {
	return itemsGenOpts(g.name+".items", g.SimpleSchema.Items)
}
func (g *simpleOpts) Required() bool {
	return g.required
}

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
	return nil, nil
}

func paramGenOpts(key string, param *spec.Parameter) (*simpleOpts, error) {
	var gopts genOpts
	if ext, ok := param.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}

	if key == "" {
		key = param.Name
	}
	return &simpleOpts{
		name:              gopts.Name,
		args:              gopts.Args,
		fieldName:         key,
		CommonValidations: param.CommonValidations,
		SimpleSchema:      param.SimpleSchema,
		required:          param.Required,
	}, nil
}

type genOpts struct {
	Name string        `mapstructure:"name"`
	Args []interface{} `mapstructure:"args"`
}

func headerGenOpts(key string, header *spec.Header) (*simpleOpts, error) {
	var gopts genOpts
	if ext, ok := header.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}
	return &simpleOpts{
		name:              gopts.Name,
		args:              gopts.Args,
		fieldName:         key,
		CommonValidations: header.CommonValidations,
		SimpleSchema:      header.SimpleSchema,
		required:          true,
	}, nil
}

func itemsGenOpts(key string, items *spec.Items) (*simpleOpts, error) {
	var gopts genOpts
	if ext, ok := items.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}
	return &simpleOpts{
		name:              gopts.Name,
		args:              gopts.Args,
		fieldName:         key,
		CommonValidations: items.CommonValidations,
		SimpleSchema:      items.SimpleSchema,
		required:          true,
	}, nil
}

// GeneratorOpts interface to capture various types that can get data generated for them.
type GeneratorOpts interface {
	// value generator name
	Name() string

	// Args for the value generator (eg. number of words in a sentence)
	// Arguments here are used to generate a valid value when no validations are specified.
	// Args are used as default setting but validations can override the args should that be necessary.
	Args() []interface{}

	// FieldName for the value generator, this is mostly used as an alternative to the name
	// for inferring which value generator to use
	FieldName() string

	// Type for the value generator to return, adids in inferring the name of the value generator
	Type() string

	// Format for the value generator to return, aids in inferring the name of the value generator
	Format() string

	// CollectionFormat how a collection is represented (csv, pipes, ...)
	CollectionFormat() string

	// Mode which kind of random data to return and to indicate which validation(s) should fail.
	// This is a bitmask so it allows for combinations of invalid values.
	Mode() StubMode

	// Maximum a numeric value can have, returns value, exclusive, defined
	Maximum() (float64, bool, bool)

	// Minimum a numeric value can have, returns value, exclusive, defined
	Minimum() (float64, bool, bool)

	// MaxLength a string can have, returns value, defined
	MaxLength() (int64, bool)

	// MinLength a string can have, returns value, defined
	MinLength() (int64, bool)

	// Pattern a string should match, returns value, defined
	Pattern() (string, bool)

	// MaxItems a collection of values can contain, returns length, defined
	MaxItems() (int64, bool)

	// MinItems a collection of values must contain, returns length, defined
	MinItems() (int64, bool)

	// UniqueItems when true the collection can't contain duplicates
	UniqueItems() bool

	// MultipleOf a numeric value should be divisible by this value, returns value, defined
	MultipleOf() (float64, bool)

	// Enum a list of acceptable values for a value, returns value, defined
	Enum() ([]interface{}, bool)

	// Items options for the members of a collection
	Items() (GeneratorOpts, error)

	// Required when true the property can't be nil
	Required() bool
}
