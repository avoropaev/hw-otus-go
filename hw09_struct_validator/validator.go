package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrConstraintIsInvalid = errors.New("constraint is invalid")
var ErrUnsupportedType = errors.New("unsupported type. Expected struct")
var ErrStringLengthIsInvalid = errors.New("string length is invalid")
var ErrStringDontMatchRegexp = errors.New("string dont match regexp")

type ErrIntLessThen int

func (e ErrIntLessThen) Error() string {
	return "value less then " + strconv.Itoa(int(e))
}

type ErrIntGreaterThen int

func (e ErrIntGreaterThen) Error() string {
	return "value greater then " + strconv.Itoa(int(e))
}

type ErrValueNotIn []string

func (e ErrValueNotIn) Error() string {
	b := strings.Builder{}
	lastI := len(e) - 1

	b.WriteString("String not in: ")

	for i, v := range e {
		b.WriteString(v)

		if i != lastI {
			b.WriteString(",")
		}
	}

	return b.String()
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	b := strings.Builder{}

	b.WriteString("validation errors:")

	for _, v := range v {
		b.WriteString("\n")
		b.WriteString(v.Field)
		b.WriteString(": ")
		b.WriteString(v.Err.Error())
	}

	return b.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	validationErrors := make(ValidationErrors, 0)

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var err error

		tag := field.Tag
		validate, ok := tag.Lookup("validate")
		if !ok {
			continue
		}

		if field.IsExported() == false {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			err = validateString(val.Field(i).String(), field.Name, validate)
		case reflect.Int:
			err = validateInt(val.Field(i).Int(), field.Name, validate)
		case reflect.Struct:
			err = Validate(val.Field(i).Interface())
			if err != nil && errors.As(err, &ValidationErrors{}) {
				validationErrors := err.(ValidationErrors)
				for i, v := range validationErrors {
					validationErrors[i].Field = field.Name + "." + v.Field
				}
			}
		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.String:
				slice, _ := val.Field(i).Interface().([]string)
				err = validateSliceOfStrings(slice, field.Name, validate)
			case reflect.Int:
				slice, _ := val.Field(i).Interface().([]int)
				err = validateSliceOfInts(slice, field.Name, validate)
			}
		}

		if err != nil {
			if errors.As(err, &ValidationErrors{}) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...)
				continue
			}

			return err
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateSliceOfStrings(slice []string, name, validate string) error {
	validationErrors := make(ValidationErrors, 0)

	for i, val := range slice {
		err := validateString(val, name+"["+strconv.Itoa(i)+"]", validate)
		if err != nil {
			if errors.As(err, &ValidationErrors{}) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...)
				continue
			}

			return err
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateSliceOfInts(slice []int, name, validate string) error {
	validationErrors := make(ValidationErrors, 0)

	for i, val := range slice {
		err := validateInt(int64(val), name+"["+strconv.Itoa(i)+"]", validate)
		if err != nil {
			if errors.As(err, &ValidationErrors{}) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...)
				continue
			}

			return err
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateString(val, name, validate string) error {
	validationErrors := make(ValidationErrors, 0)

	constraints := strings.Split(validate, "|")
	for _, constraint := range constraints {
		constraintSlice := strings.Split(constraint, ":")
		if len(constraintSlice) < 2 || len(constraintSlice[0]) == 0 || len(constraintSlice[1]) == 0 {
			return ErrConstraintIsInvalid
		}

		var validateFunc func(string, string) (error, error)

		switch constraintSlice[0] {
		case "len":
			validateFunc = validateStringLen
		case "regexp":
			validateFunc = validateStringRegexp
		case "in":
			validateFunc = validateStringIn
		default:
			continue
		}

		validationErr, unexpectedErr := validateFunc(val, constraintSlice[1])

		if unexpectedErr != nil {
			return unexpectedErr
		}

		if validationErr != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: name,
				Err:   validationErr,
			})
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateStringLen(val, constraint string) (validationErr, unexpectedErr error) {
	lengthInt, err := strconv.Atoi(constraint)
	if err != nil {
		return nil, err
	}

	if len(val) == lengthInt {
		return nil, nil
	}

	return ErrStringLengthIsInvalid, nil
}

func validateStringRegexp(val, constraint string) (validationErr, unexpectedErr error) {
	compile, err := regexp.Compile(constraint)
	if err != nil {
		return nil, err
	}

	if compile.Match([]byte(val)) {
		return nil, nil
	}

	return ErrStringDontMatchRegexp, nil
}

func validateStringIn(val, constraint string) (validationErr, unexpectedErr error) {
	listIn := strings.Split(constraint, ",")
	for _, v := range listIn {
		if val == v {
			return nil, nil
		}
	}

	return ErrValueNotIn(listIn), nil
}

func validateInt(val int64, name, validate string) error {
	validationErrors := make(ValidationErrors, 0)

	constraints := strings.Split(validate, "|")
	for _, constraint := range constraints {
		constraintSlice := strings.Split(constraint, ":")
		if len(constraintSlice) < 2 || len(constraintSlice[0]) == 0 || len(constraintSlice[1]) == 0 {
			return ErrConstraintIsInvalid
		}

		var validateFunc func(int64, string) (error, error)

		switch constraintSlice[0] {
		case "min":
			validateFunc = validateIntMin
		case "max":
			validateFunc = validateIntMax
		case "in":
			validateFunc = validateIntIn
		default:
			continue
		}

		validationErr, unexpectedErr := validateFunc(val, constraintSlice[1])

		if unexpectedErr != nil {
			return unexpectedErr
		}

		if validationErr != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: name,
				Err:   validationErr,
			})
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateIntMin(val int64, constraint string) (validationErr, unexpectedErr error) {
	min, err := strconv.Atoi(constraint)
	if err != nil {
		return nil, err
	}

	if val >= int64(min) {
		return nil, nil
	}

	return ErrIntLessThen(min), nil
}

func validateIntMax(val int64, constraint string) (validationErr, unexpectedErr error) {
	max, err := strconv.Atoi(constraint)
	if err != nil {
		return nil, err
	}

	if val <= int64(max) {
		return nil, nil
	}

	return ErrIntGreaterThen(max), nil
}

func validateIntIn(val int64, constraint string) (validationErr, unexpectedErr error) {
	listIn := strings.Split(constraint, ",")
	for _, v := range listIn {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		if val == int64(atoi) {
			return nil, nil
		}
	}

	return ErrValueNotIn(listIn), nil
}
