package stubs

import (
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/mitchellh/mapstructure"
)

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

func schemaGenOpts(key string, required bool, schema *spec.Schema) (*schemaOpts, error) {
	var gopts genOpts
	if ext, ok := schema.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}
	return &schemaOpts{
		name:      gopts.Name,
		args:      gopts.Args,
		fieldName: key,
		schema:    schema,
		required:  required,
	}, nil
}

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

type schemaOpts struct {
	schema *spec.Schema

	name string
	args []interface{}

	fieldName string
	required  bool
	mode      StubMode
}

func (s *schemaOpts) Mode() StubMode {
	return s.mode
}

func (s *schemaOpts) Args() []interface{} {
	return s.args
}

func (s *schemaOpts) Name() string {
	return s.name
}
func (s *schemaOpts) FieldName() string {
	return s.fieldName
}
func (s *schemaOpts) Maximum() (float64, bool, bool) {
	return swag.Float64Value(s.schema.Maximum), s.schema.ExclusiveMaximum, s.schema.Maximum != nil
}
func (s *schemaOpts) Minimum() (float64, bool, bool) {
	return swag.Float64Value(s.schema.Minimum), s.schema.ExclusiveMinimum, s.schema.Minimum != nil
}
func (s *schemaOpts) MaxLength() (int64, bool) {
	return swag.Int64Value(s.schema.MaxLength), s.schema.MaxLength != nil
}
func (s *schemaOpts) MinLength() (int64, bool) {
	return swag.Int64Value(s.schema.MinLength), s.schema.MinLength != nil
}
func (s *schemaOpts) Pattern() (string, bool) {
	return s.schema.Pattern, s.schema.Pattern != ""
}
func (s *schemaOpts) MaxItems() (int64, bool) {
	mx := s.schema.MaxItems
	return swag.Int64Value(mx), mx != nil
}
func (s *schemaOpts) MinItems() (int64, bool) {
	mn := s.schema.MinItems
	return swag.Int64Value(mn), mn != nil
}
func (s *schemaOpts) UniqueItems() bool {
	return s.schema.UniqueItems
}
func (s *schemaOpts) MultipleOf() (float64, bool) {
	mo := s.schema.MultipleOf
	return swag.Float64Value(mo), mo != nil
}
func (s *schemaOpts) Enum() ([]interface{}, bool) {
	enm := s.schema.Enum
	return enm, len(enm) > 0
}
func (s *schemaOpts) Type() string {
	if len(s.schema.Type) == 0 {
		return "object"
	}
	return s.schema.Type[0]
}
func (s *schemaOpts) Format() string {
	return s.schema.Format
}
func (s *schemaOpts) Items() (GeneratorOpts, error) {
	return schemaGenOpts(s.fieldName+".items", false, s.schema.Items.Schema)
}
func (s *schemaOpts) Required() bool {
	return s.required
}
