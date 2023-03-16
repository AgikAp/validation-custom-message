package validationcustommessage

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	VCM interface {
		ErrorValidationVariabel(field interface{}, rules string) error
		ErrorValidationStruct(data interface{}) (*[]ErrorRespon, bool)
		ErrorMessageOnly(err error) *[]ErrorRespon
	}

	/*
		You can see all rules validation on goplayground/validator/v10

		/f = field
		/p = param

		["uuid"] = "/f not valid uuid"
		["gte"] = "/f greater than equeal /p"
	*/
	ValidationRules map[string]string
	ErrorRespon     struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

type vcm struct {
	validationErrorRules ValidationRules
	validate             *validator.Validate
}

func New(rules ValidationRules) VCM {
	return &vcm{
		validationErrorRules: rules,
		validate:             validator.New(),
	}
}

func (v *vcm) ErrorValidationVariabel(field interface{}, rules string) error {
	err := v.validate.Var(field, rules)
	errs := v.messageGetError(strings.Split(strings.Split(err.Error(), "failed on the '")[1], "' tag")[0])

	if errs != "" {
		return errors.New(errs)
	}
	return nil
}

func (v *vcm) ErrorValidationStruct(data interface{}) (*[]ErrorRespon, bool) {
	errs := v.ErrorMessageOnly(v.validate.Struct(data))

	if len(*errs) > 0 {
		return errs, false
	}
	return errs, true
}

func (v *vcm) ErrorMessageOnly(err error) *[]ErrorRespon {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		var out []ErrorRespon
		for _, fe := range ve {
			*&out = append(*&out, ErrorRespon{
				Field:   fe.Field(),
				Message: v.messageGetError(fe.Tag(), fe.Param()),
			})
		}
		return &out
	}

	return nil
}

func (v *vcm) messageGetError(tag string, params ...string) string {
	var (
		err, param string
	)

	if len(params) > 0 {
		param = strings.Join(params, " ")
	}

	for field, msg := range v.validationErrorRules {
		if tag == field {
			r := strings.NewReplacer("/p", param)
			err = r.Replace(msg)
			break
		}
	}

	return err
}
