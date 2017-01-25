package stubs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerators_Get(t *testing.T) {
	// gen, err := newGenerator("")
	// if assert.NoError(t, err) {

	// fn, found := gen.For(&generatorOpts{
	// 	GeneratorOpts: GeneratorOpts{
	// 		Name: "characters",
	// 	},
	// })

	// if assert.True(t, found) {
	// 	res, err := fn()
	// 	if assert.NoError(t, err) {
	// 		assert.IsType(t, "", res)
	// 		assert.Len(t, res, 10)
	// 	}

	// 	res, err = fn(15)
	// 	if assert.NoError(t, err) {
	// 		assert.IsType(t, "", res)
	// 		assert.Len(t, res, 15)
	// 	}
	// }
	// }
}

func TestGeneratorsBool(t *testing.T) {
	gen, err := newGenerator("")
	if assert.NoError(t, err) {
		boolfn, found := gen.For(&simpleOpts{name: "bool"})
		if assert.True(t, found) {
			for i := 0; i < 64; i++ {
				result, err := boolfn(new(simpleOpts))
				if assert.NoError(t, err) {
					_, ok := result.(bool)
					assert.True(t, ok)
				}
			}
		}
		boolfn, found = gen.For(&simpleOpts{name: "boolean"})
		if assert.True(t, found) {
			for i := 0; i < 32; i++ {
				result, err := boolfn(new(simpleOpts))
				if assert.NoError(t, err) {
					_, ok := result.(bool)
					assert.True(t, ok)
				}
			}
		}
	}
}
