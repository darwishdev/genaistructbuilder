package internal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	genai "google.golang.org/genai"
)

type Type string

const (
	TypeString  Type = "STRING"
	TypeInteger Type = "INTEGER"
	TypeNumber  Type = "NUMBER"
	TypeBoolean Type = "BOOLEAN"
	TypeObject  Type = "OBJECT"
	TypeArray   Type = "ARRAY"
)

func BuildSchema(v any) *genai.Schema {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return buildSchemaFromType(t)
}
func buildSchemaFromType(t reflect.Type) *genai.Schema {
	s := &genai.Schema{}

	switch t.Kind() {
	case reflect.Struct:
		s.Type = genai.TypeObject
		s.Properties = map[string]*genai.Schema{}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // skip unexported
				continue
			}

			jsonTag := f.Tag.Get("json")
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName == "" {
				fieldName = f.Name
			}

			fieldSchema := buildSchemaFromType(baseType(f.Type))
			s.Properties[fieldName] = fieldSchema
			s.Required = append(s.Required, fieldName)
		}

	case reflect.Slice, reflect.Array:
		s.Type = genai.TypeArray
		s.Items = buildSchemaFromType(baseType(t.Elem()))

	case reflect.String:
		s.Type = genai.TypeString

	case reflect.Bool:
		s.Type = genai.TypeBoolean

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = genai.TypeInteger

	case reflect.Float32, reflect.Float64:
		s.Type = genai.TypeNumber

	default:
		s.Type = genai.TypeString
	}

	return s
}

func baseType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
func float32Ptr(v float32) *float32 { return &v }
func FileAdapter(
	ctx context.Context,
	fileContent []byte,
	mimeType string,
) (textPart string, mediaPart *genai.Part, err error) {
	if strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/xml" {
		return string(fileContent), nil, nil
	}
	if strings.HasPrefix(mimeType, "image/") ||
		mimeType == "application/pdf" {
		return "", &genai.Part{
			InlineData: &genai.Blob{
				Data:     fileContent,
				MIMEType: mimeType,
			},
		}, nil
	}

	// --- Case 3: Unhandled MIME Type ---
	return "", nil, fmt.Errorf("unsupported file MIME type for adapter: %s", mimeType)
}
