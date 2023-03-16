package validationcustommessage

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	VCM interface {
		ErrorValidationVariabel(field interface{}, rules string) error
		ErrorValidationStruct(data interface{}) ([]errorRespon, error)
	}

	/*
		You can see all rules validation on goplayground/validator/v10

		/f = field
		/p = param

		["uuid"] = "/f not valid uuid"
		["gte"] = "/f greater than equeal /p"
	*/
	ValidationRules map[string]string
	errorRespon     struct {
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

func (v *vcm) ErrorValidationStruct(data interface{}) ([]errorRespon, error) {
	ref := reflect.ValueOf(data)

	if ref.Kind() == reflect.Struct {
		errs := v.ErrorMessageOnly(v.validate.Struct(data))

		if len(errs) != 0 {
			findByReflect(&errs, data)
			return errs, errors.New("error field")
		}
		return nil, nil
	}
	return nil, errors.New("data not a struct")
}

func (v *vcm) ErrorMessageOnly(err error) []errorRespon {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		var out []errorRespon
		for _, fe := range ve {
			*&out = append(*&out, errorRespon{
				Field:   fe.Field(),
				Message: v.messageGetError(fe.Tag(), fe.Param()),
			})
		}
		return out
	}

	return []errorRespon{}
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

func findByReflect(data *[]errorRespon, struc interface{}) {
	ref := reflect.ValueOf(struc)
	tempData := *data

	for j, val := range tempData {
		for i := 0; i < ref.NumField(); i++ {
			if ref.Type().Field(i).Name == val.Field && ref.Type().Field(i).Tag.Get("json") != "" {
				tempData[j] = errorRespon{
					Field:   ref.Type().Field(i).Tag.Get("json"),
					Message: tempData[j].Message,
				}
			}
		}
	}
}
