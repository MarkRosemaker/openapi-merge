package merge

import (
	"fmt"
	"strings"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

// TODO: cannibalize func enrichSchema
func Schema(a, b *openapi.Schema) error {
	return merge.Schema(a, b)

	const want = `{"type":"mention":{"type":"database",`
	if strings.HasPrefix(string(b.Example), want) {
		fmt.Println("------------")
		jsonPrint("a", a)
		jsonPrint("b", b)
	}

	// if isNull(b.Example) {
	// 	// a.Nullable = true TODO
	// 	return nil
	// }

	if a.Type != b.Type {
		// if a.Type == "" {
		// 	*a = *b
		// 	return nil
		// }

		switch a.Type {
		case openapi.TypeObject:
			if isNull(a.Example) || string(a.Example) == "" {
				a.Type = b.Type
				a.Format = b.Format
				a.Example = b.Example
				// a.Nullable = true TODO
			} else if b.Type == openapi.TypeObject && string(b.Example) == "" {
				return nil
			} else {
				jsonPrint("a", a)
				fmt.Printf("a.Example: %q\n", string(a.Example))
				jsonPrint("b", b)
				fmt.Printf("b.Example: %q\n", string(b.Example))
				fmt.Println(1)
				return fmt.Errorf("schema type mismatch 1: %s != %s", a.Type, b.Type)
			}
		case openapi.TypeString:
			switch b.Type {
			case openapi.TypeObject:
				// if !isNull(b.Example) {
				// 	return fmt.Errorf("schema type mismatch 2: %s != %s", a.Type, b.Type)
				// }
			default:
				return fmt.Errorf("schema type mismatch 3: %s != %s", a.Type, b.Type)
			}
		default:
			return fmt.Errorf("schema type mismatch: %s != %s", a.Type, b.Type)
		}
	}

	if a.Format != b.Format && b.Format != "" {
		if string(a.Example) != `""` && string(b.Example) != `""` {
			if a.Format != "" {
				fmt.Printf("a.Example: %v\n", a.Example)
				fmt.Printf("b.Example: %v\n", b.Example)
				fmt.Println(5)

				return fmt.Errorf("schema format mismatch: %q != %q", a.Format, b.Format)
			} // else: don't overwrite empty format if a has an example,
			// i.e. does not match b's format
		} else if a.Format == "" {
			a.Format = b.Format
			a.Example = b.Example
		}
	}

	// AllOf SchemaRefList `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	// Min *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// Max *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// Pattern *regexp.Regexp `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	// Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	// MinItems uint `json:"minItems,omitzero" yaml:"minItems,omitempty"`
	// MaxItems *uint `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// Properties SchemaRefs `json:"properties,omitzero" yaml:"properties,omitempty"`
	// Required             []string   `json:"required,omitempty" yaml:"required,omitempty"`
	// AdditionalProperties *SchemaRef `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	// ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
	// ContentEncoding  string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`
	// Default any `json:"default,omitempty" yaml:"default,omitempty"`

	return nil
}
