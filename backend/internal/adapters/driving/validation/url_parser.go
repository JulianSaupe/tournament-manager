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

		// Parse path parameters
		if pathTag := fieldType.Tag.Get("path"); pathTag != "" {
			if value := chi.URLParam(r, pathTag); value != "" {
				if err := p.setFieldValue(field, value); err != nil {
					return fmt.Errorf("failed to set path param %s: %w", pathTag, err)
				}
			}
		}

		// Parse query parameters
		if queryTag := fieldType.Tag.Get("query"); queryTag != "" {
			if value := r.URL.Query().Get(queryTag); value != "" {
				if err := p.setFieldValue(field, value); err != nil {
					return fmt.Errorf("failed to set query param %s: %w", queryTag, err)
				}
			}
		}

		// Set default values
		if defaultTag := fieldType.Tag.Get("default"); defaultTag != "" && p.isZeroValue(field) {
			if err := p.setFieldValue(field, defaultTag); err != nil {
				return fmt.Errorf("failed to set default value for %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

func (p *URLParser) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Slice:
		// Handle comma-separated values
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(strings.Split(value, ",")))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func (p *URLParser) isZeroValue(field reflect.Value) bool {
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
