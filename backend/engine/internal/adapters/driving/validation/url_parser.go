package validation

import (
	"Tournament/internal/domain"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	pathTag    = "path"
	queryTag   = "query"
	defaultTag = "default"
)

// URLParser handles parsing URL parameters into structs
type URLParser struct {
	validator *validator.Validate
}

func NewURLParser() *URLParser {
	return &URLParser{
		validator: validator.New(),
	}
}

func (p *URLParser) parseURLParams(r *http.Request, params interface{}) error {
	if err := p.parseParameters(r, params); err != nil {
		return err
	}
	if err := p.validator.Struct(params); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

func (p *URLParser) parseParameters(r *http.Request, params interface{}) error {
	v := reflect.ValueOf(params).Elem()
	t := reflect.TypeOf(params).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		if err := p.processField(r, field, fieldType); err != nil {
			return err
		}
	}
	return nil
}

func (p *URLParser) processField(r *http.Request, field reflect.Value, fieldType reflect.StructField) error {
	if err := p.processPathParameter(r, field, fieldType); err != nil {
		return err
	}

	if err := p.processQueryParameter(r, field, fieldType); err != nil {
		return err
	}

	if err := p.processDefaultValue(field, fieldType); err != nil {
		return err
	}

	return nil
}

func (p *URLParser) processPathParameter(r *http.Request, field reflect.Value, fieldType reflect.StructField) error {
	pathTagValue := fieldType.Tag.Get(pathTag)
	if pathTagValue == "" {
		return nil
	}

	value := chi.URLParam(r, pathTagValue)
	if value == "" {
		return nil
	}

	if err := p.setFieldValue(field, value); err != nil {
		return fmt.Errorf("failed to set path param %s: %w", pathTagValue, err)
	}
	return nil
}

func (p *URLParser) processQueryParameter(r *http.Request, field reflect.Value, fieldType reflect.StructField) error {
	queryTagValue := fieldType.Tag.Get(queryTag)
	if queryTagValue == "" {
		return nil
	}

	value := r.URL.Query().Get(queryTagValue)
	if value == "" {
		return nil
	}

	if err := p.setFieldValue(field, value); err != nil {
		return fmt.Errorf("failed to set query param %s: %w", queryTagValue, err)
	}
	return nil
}

func (p *URLParser) processDefaultValue(field reflect.Value, fieldType reflect.StructField) error {
	defaultValue := fieldType.Tag.Get(defaultTag)
	if defaultValue == "" || !p.isFieldEmpty(field) {
		return nil
	}

	if err := p.setFieldValue(field, defaultValue); err != nil {
		return fmt.Errorf("failed to set default value for %s: %w", fieldType.Name, err)
	}
	return nil
}

func (p *URLParser) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		return p.setStringValue(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return p.setIntValue(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return p.setUintValue(field, value)
	case reflect.Bool:
		return p.setBoolValue(field, value)
	case reflect.Float32, reflect.Float64:
		return p.setFloatValue(field, value)
	case reflect.Slice:
		return p.setSliceValue(field, value)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
}

func (p *URLParser) setStringValue(field reflect.Value, value string) error {
	field.SetString(value)
	return nil
}

func (p *URLParser) setIntValue(field reflect.Value, value string) error {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	field.SetInt(intValue)
	return nil
}

func (p *URLParser) setUintValue(field reflect.Value, value string) error {
	uintValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(uintValue)
	return nil
}

func (p *URLParser) setBoolValue(field reflect.Value, value string) error {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(boolValue)
	return nil
}

func (p *URLParser) setFloatValue(field reflect.Value, value string) error {
	floatValue, err := strconv.ParseFloat(value, field.Type().Bits())
	if err != nil {
		return err
	}
	field.SetFloat(floatValue)
	return nil
}

func (p *URLParser) setSliceValue(field reflect.Value, value string) error {
	// Handle comma-separated values
	if field.Type().Elem().Kind() == reflect.String {
		field.Set(reflect.ValueOf(strings.Split(value, ",")))
		return nil
	}
	return fmt.Errorf("unsupported slice element type: %s", field.Type().Elem().Kind())
}

func (p *URLParser) isFieldEmpty(field reflect.Value) bool {
	return field.IsZero()
}

var DefaultURLParser = NewURLParser()

func ParseURLParams[T any](r *http.Request) (*T, error) {
	var params T
	if err := DefaultURLParser.parseURLParams(r, &params); err != nil {
		return nil, err
	}
	return &params, nil
}

func ValidateURLParams[T any](r *http.Request) *T {
	params, err := ParseURLParams[T](r)
	if err != nil {
		panic(domain.NewInvalidParameterError(err.Error()))
	}
	return params
}
