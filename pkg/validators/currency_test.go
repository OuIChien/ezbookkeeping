package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidCurrency(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validCurrency", ValidCurrency)
	assert.Nil(t, err)

	err = validate.Var("CNY", "validCurrency")
	assert.Nil(t, err)

	err = validate.Var("USD", "validCurrency")
	assert.Nil(t, err)

	err = validate.Var("---", "validCurrency")
	assert.Nil(t, err)

	err = validate.Var("XXX", "validCurrency")
	assert.Nil(t, err)

	err = validate.Var("RMB", "validCurrency")
	assert.Nil(t, err)
}

func TestInvalidCurrency(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validCurrency", ValidCurrency)
	assert.Nil(t, err)

	err = validate.Var("US$", "validCurrency")
	assert.NotNil(t, err)

	err = validate.Var("TOOLONGNAME", "validCurrency")
	assert.NotNil(t, err)

	err = validate.Var("", "validCurrency")
	assert.NotNil(t, err)

	err = validate.Var("-", "validCurrency")
	assert.NotNil(t, err)
}
