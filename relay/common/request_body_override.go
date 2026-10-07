package common

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const paramOverrideContextRequestBody = "__request_body_override"

// RequestBodyOverrideOptions controls encoding after parameter overrides.
// It is relay metadata populated by channel operations, never an upstream field.
type RequestBodyOverrideOptions struct {
	Format     string            `json:"format"`
	FieldTypes map[string]string `json:"field_types,omitempty"`
}

func parseRequestBodyOverrideOptions(value interface{}) (*RequestBodyOverrideOptions, error) {
	data, err := common.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("set_request_body: invalid options: %w", err)
	}
	if !gjson.ParseBytes(data).IsObject() {
		return nil, fmt.Errorf("set_request_body: value must be a JSON object")
	}
	var options RequestBodyOverrideOptions
	if err := common.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("set_request_body: invalid options: %w", err)
	}
	options.Format = strings.ToLower(strings.TrimSpace(options.Format))
	switch options.Format {
	case "json", "multipart":
	default:
		return nil, fmt.Errorf("set_request_body: format must be json or multipart")
	}
	for path, fieldType := range options.FieldTypes {
		if strings.TrimSpace(path) == "" || path != strings.TrimSpace(path) || strings.ContainsAny(path, "*?#|") {
			return nil, fmt.Errorf("set_request_body: invalid field path %q", path)
		}
		switch fieldType {
		case "integer", "number", "boolean", "string":
		default:
			return nil, fmt.Errorf("set_request_body: unsupported field type %q for %s", fieldType, path)
		}
	}
	return &options, nil
}

// BuildRequestBodyOverride preserves all overridden fields, including explicit
// zero/false/null values. Deleted fields are not reconstructed from task DTOs.
// JSON source bodies are required; multipart here contains text fields only.
func BuildRequestBodyOverride(data []byte, upstreamModel string, options *RequestBodyOverrideOptions) (io.Reader, string, error) {
	if options == nil {
		return nil, "", fmt.Errorf("request body override options are required")
	}
	if !gjson.ValidBytes(data) || !gjson.ParseBytes(data).IsObject() {
		return nil, "", fmt.Errorf("request body override requires a JSON object")
	}
	var err error
	if upstreamModel != "" {
		data, err = sjson.SetBytes(data, "model", upstreamModel)
		if err != nil {
			return nil, "", err
		}
	}
	data, err = convertRequestBodyFieldTypes(data, options.FieldTypes)
	if err != nil {
		return nil, "", err
	}
	switch options.Format {
	case "json":
		return bytes.NewReader(data), "application/json", nil
	case "multipart":
		return encodeMultipartRequestBody(data)
	default:
		return nil, "", fmt.Errorf("unsupported request body format %q", options.Format)
	}
}

func convertRequestBodyFieldTypes(data []byte, fieldTypes map[string]string) ([]byte, error) {
	paths := make([]string, 0, len(fieldTypes))
	for path := range fieldTypes {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		value := gjson.GetBytes(data, path)
		if !value.Exists() || value.Type == gjson.Null {
			continue
		}
		converted, err := convertRequestBodyScalar(value, fieldTypes[path])
		if err != nil {
			return nil, fmt.Errorf("convert request field %s to %s: %w", path, fieldTypes[path], err)
		}
		data, err = sjson.SetBytes(data, path, converted)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func convertRequestBodyScalar(value gjson.Result, fieldType string) (interface{}, error) {
	if value.IsObject() || value.IsArray() {
		return nil, fmt.Errorf("value must be a scalar")
	}
	text := value.String()
	if value.Type != gjson.String {
		text = value.Raw
	}
	switch fieldType {
	case "string":
		return text, nil
	case "integer":
		return strconv.ParseInt(strings.TrimSpace(text), 10, 64)
	case "number":
		number, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			return nil, err
		}
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, fmt.Errorf("value must be a finite number")
		}
		return number, nil
	case "boolean":
		return strconv.ParseBool(strings.TrimSpace(text))
	default:
		return nil, fmt.Errorf("unsupported field type %q", fieldType)
	}
}

func encodeMultipartRequestBody(data []byte) (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	var writeErr error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		// Strings are plain text; arrays/objects and other scalars use JSON text.
		text := value.Raw
		if value.Type == gjson.String {
			text = value.String()
		}
		writeErr = writer.WriteField(key.String(), text)
		return writeErr == nil
	})
	if writeErr != nil {
		_ = writer.Close()
		return nil, "", writeErr
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return &buf, writer.FormDataContentType(), nil
}
