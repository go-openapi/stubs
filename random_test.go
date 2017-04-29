package stubs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerators_Characters(t *testing.T) {
	gen, err := newGenerator("")
	if assert.NoError(t, err) {
		opts := &simpleOpts{name: "characters"}
		fn, found := gen.For(opts)

		if assert.True(t, found) {
			res, err := fn(opts)
			if assert.NoError(t, err) {
				assert.IsType(t, "", res)
				assert.Len(t, res, 10)
			}
			opts.args = append(opts.args, 15)
			res, err = fn(opts)
			if assert.NoError(t, err) {
				assert.IsType(t, "", res)
				assert.Len(t, res, 15)
			}
		}
	}
}

func TestGeneratorsBool(t *testing.T) {
	gen, err := newGenerator("")
	if assert.NoError(t, err) {
		boolfn, found := gen.For(&simpleOpts{name: "bool"})
		if assert.True(t, found) {
			for i := 0; i < 64; i++ {
				result, err := boolfn(new(simpleOpts))
				if assert.NoError(t, err) {
					assert.IsType(t, true, result)
				}
			}
		}
		boolfn, found = gen.For(&simpleOpts{name: "boolean"})
		if assert.True(t, found) {
			for i := 0; i < 32; i++ {
				result, err := boolfn(new(simpleOpts))
				if assert.NoError(t, err) {
					assert.IsType(t, true, result)
				}
			}
		}
	}
}
