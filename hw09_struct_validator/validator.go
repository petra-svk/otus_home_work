package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Error: %v\n", v.Field, v.Err)
}

func (v *ValidationError) Unwrap() error {
	return v.Err
}

var (
	ErrNotStruct          = errors.New("wrong type struct")
	ErrRegexpCompile      = errors.New("wrong regexp")
	ErrConvertionStrToInt = errors.New("error convert string to int")
	ErrValidMethod        = errors.New("validation method is not supported")

	ErrLen     = errors.New("wrong len")
	ErrRegexp  = errors.New("value is not match to regexp")
	ErrInclude = errors.New("wrong value in enumeration")
	ErrMin     = errors.New("value is less then min")
	ErrMax     = errors.New("value is more then max")
)

type ValueData struct {
	FieldName string
	Tag       string
	Value     interface{}
}

type ValidationErrors []ValidationError

func (v *ValidationErrors) Error() string {
	var b strings.Builder
	for _, err := range *v {
		b.WriteString(fmt.Sprintf("Field: %s, Error: %v\n", err.Field, err.Err))
	}
	return b.String()
}

func checkLen(field, v, length string) error {
	i, err := strconv.Atoi(length)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if len(v) != i {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual length %d", ErrLen, len(v)),
		}
	}
	return nil
}

func checkRegexp(field, v, r string) error {
	re, err := regexp.Compile(r)
	if err != nil {
		return ErrRegexpCompile
	}
	if !re.MatchString(v) {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%s, regexp=`%s`", ErrRegexp, v, r),
		}
	}
	return nil
}

func checkIn(field string, v interface{}, inc string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.String {
		s := strings.Split(inc, ",")
		dict := make(map[string]struct{}, len(s))
		for _, val := range s {
			dict[val] = struct{}{}
		}
		if _, ok := dict[rv.String()]; !ok {
			return &ValidationError{
				Field: field,
				Err:   fmt.Errorf("%w - value=%s, enum=%s", ErrInclude, rv.String(), inc),
			}
		}
	}
	if rv.Kind() == reflect.Int {
		s := strings.Split(inc, ",")
		dict := make(map[int]struct{}, len(s))
		for _, val := range s {
			i, err := strconv.Atoi(val)
			if err != nil {
				return ErrConvertionStrToInt
			}
			dict[i] = struct{}{}
		}
		if _, ok := dict[int(rv.Int())]; !ok {
			return &ValidationError{
				Field: field,
				Err:   fmt.Errorf("%w - value=%d, enum=%s", ErrInclude, rv.Int(), inc),
			}
		}
	}
	return nil
}

func minVal(field string, v int, min string) error {
	i, err := strconv.Atoi(min)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if v < i {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%d, min=%d", ErrMin, v, i),
		}
	}
	return nil
}

func maxVal(field string, v int, max string) error {
	i, err := strconv.Atoi(max)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if v > i {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%d, max=%d", ErrMax, v, i),
		}
	}
	return nil
}

func checkVal(field, tag string, v interface{}) error {
	// our check methods
	vMap := map[string]interface{}{
		"in":     checkIn,
		"len":    checkLen,
		"regexp": checkRegexp,
		"min":    minVal,
		"max":    maxVal,
	}
	var vErr ValidationErrors
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		var err error
		r := strings.Split(rule, ":")
		switch r[0] {
		case "in":
			err = vMap["in"].(func(string, interface{}, string) error)(field, v, r[1])
		case "len":
			err = vMap["len"].(func(string, string, string) error)(field, v.(string), r[1])
		case "regexp":
			err = vMap["regexp"].(func(string, string, string) error)(field, v.(string), r[1])
		case "min":
			err = vMap["min"].(func(string, int, string) error)(field, v.(int), r[1])
		case "max":
			err = vMap["max"].(func(string, int, string) error)(field, v.(int), r[1])
		default:
			return ErrValidMethod
		}
		var e *ValidationError
		if err != nil {
			if errors.As(err, &e) {
				vErr = append(vErr, *e)
			} else {
				return err
			}
		}
	}
	return &vErr
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var validErrors ValidationErrors // hub for validation errors
	t := rv.Type()
	var values []ValueData
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)     // reflect.StructField
		fv := rv.Field(i)       // reflect.Value
		if !fv.CanInterface() { // field is private
			continue
		}
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		switch fv.Type().Kind().String() {
		case "int":
			values = append(values, ValueData{field.Name, tag, int(fv.Int())})
		case "string":
			values = append(values, ValueData{field.Name, tag, fv.String()})
		case "slice":
			if fv.Type().String() == "[]string" {
				strSlice := fv.Interface().([]string)
				for i := 0; i < len(strSlice); i++ {
					values = append(values, ValueData{field.Name, tag, strSlice[i]})
				}
			}
			if fv.Type().String() == "[]int" {
				intSlice := fv.Interface().([]int)
				for i := 0; i < len(intSlice); i++ {
					values = append(values, ValueData{field.Name, tag, intSlice[i]})
				}
			}
		}
	}
	for _, val := range values {
		err := checkVal(val.FieldName, val.Tag, val.Value)
		if err != nil {
			var e *ValidationErrors
			if errors.As(err, &e) {
				validErrors = append(validErrors, *e...)
			} else {
				return err
			}
		}
	}

	if len(validErrors) > 0 {
		return &validErrors
	}
	return nil
}
