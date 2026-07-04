package merge_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestMediaType(t *testing.T) {
	t.Parallel()

	a := &openapi.MediaType{}
	b := &openapi.MediaType{
		Schema:  &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
		Example: jsontext.Value(`"foo"`),
	}

	if err := merge.MediaType(a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// the result must land in a, the first argument
	if got := a.Schema.Value.Type; got != openapi.TypeString {
		t.Fatalf("expected type %q, got %q", openapi.TypeString, got)
	}

	if string(a.Example) != `"foo"` {
		t.Fatalf("expected example %q, got %q", `"foo"`, string(a.Example))
	}
}
