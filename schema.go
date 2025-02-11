package merge

import (
	"fmt"
	"slices"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func Schema(a, b *openapi.Schema) error {
	a.Title = mergeString(a.Title, b.Title)
	a.Description = mergeString(a.Description, b.Description)

	if a.Example == nil {
		a.Example = b.Example
	}

	if isNull(b.Example) && b.AdditionalProperties == nil &&
		b.Properties == nil {
		// else: b is a string or other type that got null'ed
		b.Type = a.Type
		b.Format = a.Format
		// TODO: a.Nullable = true
	}

	if a.Type != b.Type {
		jsonPrint("a", a)
		jsonPrint("b", b)

		return &errpath.ErrField{Field: "type", Err: fmt.Errorf("%q != %q", a.Type, b.Type)}
	}

	if a.Format != b.Format {
		jsonPrint("a", a)
		jsonPrint("b", b)

		return &errpath.ErrField{Field: "format", Err: fmt.Errorf("%q != %q", a.Format, b.Format)}
	}

	switch a.Type {
	case openapi.TypeString:
		// add the example from b to the enums of a
		if a.Enum != nil && b.Example != nil {
			ex := ""
			if err := json.Unmarshal(b.Example, &ex); err != nil {
				return &errpath.ErrField{Field: "example", Err: err}
			}

			if !slices.Contains(a.Enum, ex) {
				a.Enum = append(a.Enum, ex)
			}
		}
	case openapi.TypeObject:
		if a.AdditionalProperties != nil {
			// is a string map
			if b.AdditionalProperties != nil {
				if err := Schema(a.AdditionalProperties.Value, b.AdditionalProperties.Value); err != nil {
					return &errpath.ErrField{Field: "additionalProperties", Err: err}
				}
			}

			// merge all property values with prop
			for _, prop := range b.Properties {
				if err := Schema(a.AdditionalProperties.Value, prop.Value); err != nil {
					return &errpath.ErrField{Field: "additionalProperties", Err: err}
				}
			}
		} else if a.Properties == nil {
			a.Properties = b.Properties // simply set the properties
		} else {
			for k, propB := range b.Properties.ByIndex() {
				propA, ok := a.Properties[k]
				if !ok {
					a.Properties.Set(k, propB) // add the property
					continue
				}

				// merge the properties
				if err := Schema(propA.Value, propB.Value); err != nil {
					return err
				}
			}
		}
	case openapi.TypeBoolean: // nothing to do
	case openapi.TypeArray:
		if err := Schema(a.Items.Value, b.Items.Value); err != nil {
			return err
		}
	default:
		jsonPrint("a", a)
		jsonPrint("b", b)

		return &errpath.ErrField{Field: "type", Err: fmt.Errorf("%q unimplemented", a.Type)}
	}

	if len(a.AllOf) > 0 {
		return &errpath.ErrField{Field: "allOf", Err: fmt.Errorf("not implemented")}
	}

	// // validate if format is valid for type
	// switch s.Format {
	// case "": // no format
	// case FormatInt32, FormatInt64:
	// 	if s.Type != TypeInteger {
	// 		return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
	// 			Value:   s.Format,
	// 			Message: fmt.Sprintf("only valid for integer type, got %s", s.Type),
	// 		}}
	// 	}
	// case FormatFloat, FormatDouble:
	// 	if s.Type != TypeNumber {
	// 		return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
	// 			Value:   s.Format,
	// 			Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
	// 		}}
	// 	}
	// case FormatDateTime, FormatEmail, FormatPassword,
	// 	FormatUUID, FormatURI, FormatURIRef, FormatZipCode,
	// 	FormatIPv4, FormatIPv6:
	// 	if s.Type != TypeString {
	// 		return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
	// 			Value:   s.Format,
	// 			Message: fmt.Sprintf("only valid for string type, got %s", s.Type),
	// 		}}
	// 	}
	// case FormatDuration:
	// 	switch s.Type {
	// 	case TypeInteger, TypeString:
	// 	default:
	// 		return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
	// 			Value:   s.Format,
	// 			Message: fmt.Sprintf("only valid for integer or string type, got %s", s.Type),
	// 		}}
	// 	}
	// default:
	// 	return fmt.Errorf("unimplemented format: %s", s.Format)
	// }

	// for i, v := range s.AllOf {
	// 	if err := v.Validate(); err != nil {
	// 		return &errpath.ErrField{
	// 			Field: "allOf",
	// 			Err:   &errpath.ErrIndex{Index: i, Err: err},
	// 		}
	// 	}
	// }

	// // Integer / Number

	// // validate min and max
	// if s.Type == TypeInteger {
	// 	if s.Min != nil && *s.Min != float64(int(*s.Min)) {
	// 		return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
	// 			Value:   *s.Min,
	// 			Message: "not an integer",
	// 		}}
	// 	}

	// 	if s.Max != nil && *s.Max != float64(int(*s.Max)) {
	// 		return &errpath.ErrField{Field: "maximum", Err: &errpath.ErrInvalid[float64]{
	// 			Value:   *s.Max,
	// 			Message: "not an integer",
	// 		}}
	// 	}
	// }

	// if s.Type == TypeNumber || s.Type == TypeInteger {
	// 	if s.Min != nil && s.Max != nil && *s.Min > *s.Max {
	// 		return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
	// 			Value:   *s.Min,
	// 			Message: fmt.Sprintf("minimum is greater than maximum (%v > %v)", *s.Min, *s.Max),
	// 		}}
	// 	}
	// } else if s.Min != nil {
	// 	return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
	// 		Value:   *s.Min,
	// 		Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
	// 	}}
	// } else if s.Max != nil {
	// 	return &errpath.ErrField{Field: "maximum", Err: &errpath.ErrInvalid[float64]{
	// 		Value:   *s.Max,
	// 		Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
	// 	}}
	// }

	// // String

	// if s.Type != TypeString && s.Enum != nil {
	// 	return &errpath.ErrField{Field: "enum", Err: &errpath.ErrInvalid[string]{
	// 		Message: fmt.Sprintf("only valid for string type, got %s", s.Type),
	// 	}}
	// }

	// // Array

	// // validate min and max items
	// if s.Type == TypeArray {
	// 	if s.MaxItems != nil && s.MinItems > *s.MaxItems {
	// 		return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
	// 			Value:   s.MinItems,
	// 			Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
	// 		}}
	// 	}

	// 	if s.Items == nil {
	// 		return &errpath.ErrField{Field: "items", Err: &errpath.ErrRequired{}}
	// 	}

	// 	// empty schema for items indicates a media type of application/octet-stream.
	// 	if !s.Items.Value.isEmpty() {
	// 		if err := s.Items.Validate(); err != nil {
	// 			return &errpath.ErrField{Field: "items", Err: err}
	// 		}
	// 	}
	// } else if s.MinItems != 0 {
	// 	return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
	// 		Value:   s.MinItems,
	// 		Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
	// 	}}
	// } else if s.MaxItems != nil {
	// 	return &errpath.ErrField{Field: "maxItems", Err: &errpath.ErrInvalid[uint]{
	// 		Value:   *s.MaxItems,
	// 		Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
	// 	}}
	// } else if s.Items != nil {
	// 	return &errpath.ErrField{Field: "items", Err: &errpath.ErrInvalid[string]{
	// 		Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
	// 	}}
	// }

	// // Object

	// if s.Type == TypeObject {
	// 	if err := s.Properties.Validate(); err != nil {
	// 		return &errpath.ErrField{Field: "properties", Err: err}
	// 	}

	// 	for i, r := range s.Required {
	// 		if _, ok := s.Properties[r]; ok {
	// 			continue
	// 		}

	// 		return &errpath.ErrField{
	// 			Field: "required",
	// 			Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[string]{
	// 				Value:   r,
	// 				Message: "property does not exist",
	// 			}},
	// 		}
	// 	}

	// 	if s.AdditionalProperties != nil {
	// 		if err := s.AdditionalProperties.Validate(); err != nil {
	// 			return &errpath.ErrField{Field: "additionalProperties", Err: err}
	// 		}
	// 	}
	// } else if s.Properties != nil {
	// 	return &errpath.ErrField{Field: "properties", Err: &errpath.ErrInvalid[string]{
	// 		Message: fmt.Sprintf("only valid for object type, got %s", s.Type),
	// 	}}
	// } else if s.AdditionalProperties != nil {
	// 	return &errpath.ErrField{Field: "additionalProperties", Err: &errpath.ErrInvalid[string]{
	// 		Message: fmt.Sprintf("only valid for object type, got %s", s.Type),
	// 	}}
	// }

	// // validate default
	// switch dflt := s.Default.(type) {
	// case nil: // empty
	// case string:
	// 	if s.Type != TypeString {
	// 		return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[string]{
	// 			Value:   dflt,
	// 			Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
	// 		}}
	// 	}

	// 	if s.Enum != nil {
	// 		if !slices.Contains(s.Enum, dflt) {
	// 			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[string]{
	// 				Value:   dflt,
	// 				Message: fmt.Sprintf("is not one of the enums (%q)", s.Enum),
	// 			}}
	// 		}
	// 	}
	// case float64:
	// 	switch s.Type {
	// 	case TypeNumber: // fits
	// 	case TypeInteger:
	// 		if asInt := int(dflt); dflt != float64(asInt) {
	// 			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[float64]{
	// 				Value:   dflt,
	// 				Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
	// 			}}
	// 		} else {
	// 			s.Default = asInt // set to int version
	// 		}
	// 	default:
	// 		return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[float64]{
	// 			Value:   dflt,
	// 			Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
	// 		}}
	// 	}
	// case int:
	// 	switch s.Type {
	// 	case TypeNumber, TypeInteger: // fits
	// 	default:
	// 		return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[int]{
	// 			Value:   dflt,
	// 			Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
	// 		}}
	// 	}
	// default:
	// 	return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
	// 		Value:   s.Default,
	// 		Message: fmt.Sprintf("unknown type %T", s.Default),
	// 	}}
	// }

	return nil
}

func jsonPrint(name string, s *openapi.Schema) {
	data, _ := json.Marshal(s)
	fmt.Printf("%s: %s\n", name, string(data))
}

func isNull(j jsontext.Value) bool {
	return string(j) == jsontext.Null.String()
}

// func (l *loader) resolveSchema(s *Schema) error {
// 	if err := l.resolveSchemaRefList(s.AllOf); err != nil {
// 		return &errpath.ErrField{Field: "allOf", Err: err}
// 	}

// 	if s.Items != nil {
// 		if err := l.resolveSchemaRef(s.Items); err != nil {
// 			return &errpath.ErrField{Field: "items", Err: err}
// 		}
// 	}

// 	if err := l.resolveSchemaRefs(s.Properties); err != nil {
// 		return &errpath.ErrField{Field: "properties", Err: err}
// 	}

// 	if s.AdditionalProperties != nil {
// 		if err := l.resolveSchemaRef(s.AdditionalProperties); err != nil {
// 			return &errpath.ErrField{Field: "additionalProperties", Err: err}
// 		}
// 	}

// 	return nil
// }
