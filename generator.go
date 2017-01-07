package stubs

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/mitchellh/mapstructure"
)

// StubMode for generating data
type StubMode uint64

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

// type StubOpts struct {
// 	Mode       StubMode
// 	Language   string
// 	Descriptor interface{} // spec.Schema, spec.Header, spec.Parameter
// 	Target     interface{}
// }

type generatorOpts struct {
	spec.CommonValidations
	spec.SimpleSchema

	name      string
	args      []interface{}
	fieldName string
	required  bool
}

func (g *generatorOpts) Args() []interface{} {
	return g.args
}

func (g *generatorOpts) GoType() reflect.Type {
	return nil
}
func (g *generatorOpts) CollectionFormat() string {
	return g.SimpleSchema.CollectionFormat
}

func (g *generatorOpts) Name() string {
	return g.name
}
func (g *generatorOpts) FieldName() string {
	return g.fieldName
}
func (g *generatorOpts) Maximum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Maximum), g.CommonValidations.ExclusiveMaximum, g.CommonValidations.Maximum != nil
}
func (g *generatorOpts) Minimum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Minimum), g.CommonValidations.ExclusiveMinimum, g.CommonValidations.Minimum != nil
}
func (g *generatorOpts) MaxLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MaxLength), g.CommonValidations.MaxLength != nil
}
func (g *generatorOpts) MinLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MinLength), g.CommonValidations.MinLength != nil
}
func (g *generatorOpts) Pattern() (string, bool) {
	return g.CommonValidations.Pattern, g.CommonValidations.Pattern != ""
}
func (g *generatorOpts) MaxItems() (int64, bool) {
	mx := g.CommonValidations.MaxItems
	return swag.Int64Value(mx), mx != nil
}
func (g *generatorOpts) MinItems() (int64, bool) {
	mn := g.CommonValidations.MinItems
	return swag.Int64Value(mn), mn != nil
}
func (g *generatorOpts) UniqueItems() bool {
	return g.CommonValidations.UniqueItems
}
func (g *generatorOpts) MultipleOf() (float64, bool) {
	mo := g.CommonValidations.MultipleOf
	return swag.Float64Value(mo), mo != nil
}
func (g *generatorOpts) Enum() ([]interface{}, bool) {
	enm := g.CommonValidations.Enum
	return enm, len(enm) > 0
}
func (g *generatorOpts) Type() string {
	return g.SimpleSchema.Type
}
func (g *generatorOpts) Format() string {
	return g.SimpleSchema.Format
}
func (g *generatorOpts) Items() (GeneratorOpts, error) {
	return itemsGenOpts(g.name+".items", g.SimpleSchema.Items)
}
func (g *generatorOpts) Required() bool {
	return g.required
}

type Stubbing struct {
	Language string
}

// Generate a stub into the opts.Target
func (s *Stubbing) Generate(key string, descriptor interface{}) (interface{}, error) {

	// if err := s.checkValid(opts); err != nil {
	// 	return err
	// }

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
		return nil, fmt.Errorf("%T is unsupported for stubbing", descriptor)
	}
}

func (s *Stubbing) GenParameter(key string, param *spec.Parameter) (interface{}, error) {
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

	return datagen(gopts.Args()...)
}

func (s *Stubbing) GenHeader(key string, header *spec.Header) (interface{}, error) {
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

	return datagen(gopts.Args()...)
}

func (s *Stubbing) GenSchema(key string, schema *spec.Schema) (interface{}, error) {
	return nil, nil
}

func paramGenOpts(key string, param *spec.Parameter) (*generatorOpts, error) {
	var gopts genOpts
	if ext, ok := param.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}

	if key == "" {
		key = param.Name
	}
	return &generatorOpts{
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

func headerGenOpts(key string, header *spec.Header) (*generatorOpts, error) {
	var gopts genOpts
	if ext, ok := header.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}
	return &generatorOpts{
		name:              gopts.Name,
		args:              gopts.Args,
		fieldName:         key,
		CommonValidations: header.CommonValidations,
		SimpleSchema:      header.SimpleSchema,
		required:          true,
	}, nil
}

func itemsGenOpts(key string, items *spec.Items) (*generatorOpts, error) {
	var gopts genOpts
	if ext, ok := items.Extensions["x-datagen"]; ok {
		if err := mapstructure.WeakDecode(ext, &gopts); err != nil {
			return nil, err
		}
	}
	return &generatorOpts{
		name:              gopts.Name,
		args:              gopts.Args,
		fieldName:         key,
		CommonValidations: items.CommonValidations,
		SimpleSchema:      items.SimpleSchema,
		required:          true,
	}, nil
}

type GeneratorOpts interface {
	Name() string
	Args() []interface{}
	FieldName() string
	Maximum() (float64, bool, bool)
	Minimum() (float64, bool, bool)
	MaxLength() (int64, bool)
	MinLength() (int64, bool)
	Pattern() (string, bool)
	MaxItems() (int64, bool)
	MinItems() (int64, bool)
	UniqueItems() bool
	MultipleOf() (float64, bool)
	Enum() ([]interface{}, bool)
	Type() string
	Format() string
	Items() (GeneratorOpts, error)
	CollectionFormat() string
	Required() bool
}
