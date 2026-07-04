package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestContent(t *testing.T) {
	t.Parallel()

	// a starts out nil: merging must still make the new entry visible to the
	// caller, even though maps are normally passed by value.
	var a openapi.Content

	mt := &openapi.MediaType{Schema: &openapi.SchemaRef{Value: &openapi.Schema{
		Type: openapi.TypeString,
	}}}
	b := openapi.Content{"application/json": mt}

	if err := merge.Content(&a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a))
	}

	if a["application/json"] != mt {
		t.Fatalf("expected the same media type pointer to be present in a")
	}
}

func TestContent_MergeExistingKey(t *testing.T) {
	t.Parallel()

	a := openapi.Content{"application/json": &openapi.MediaType{
		Schema: &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
	}}
	b := openapi.Content{"application/json": &openapi.MediaType{
		Schema: &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
	}}

	if err := merge.Content(&a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a))
	}

	if got := a["application/json"].Schema.Value.Type; got != openapi.TypeString {
		t.Fatalf("expected type %q, got %q", openapi.TypeString, got)
	}
}
