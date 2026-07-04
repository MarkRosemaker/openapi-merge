package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestSchemaRefs(t *testing.T) {
	t.Parallel()

	// a starts out nil: merging must still make the new entry visible to the
	// caller, even though maps are normally passed by value.
	var a openapi.SchemaRefs

	sr := &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}}
	b := openapi.SchemaRefs{"name": sr}

	if err := merge.SchemaRefs(&a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a))
	}

	if a["name"] != sr {
		t.Fatalf("expected the same schema ref pointer to be present in a")
	}
}

func TestSchemaRefs_MergeExistingKey(t *testing.T) {
	t.Parallel()

	a := openapi.SchemaRefs{"name": &openapi.SchemaRef{
		Value: &openapi.Schema{Type: openapi.TypeString},
	}}
	b := openapi.SchemaRefs{"name": &openapi.SchemaRef{
		Value: &openapi.Schema{Type: openapi.TypeString},
	}}

	if err := merge.SchemaRefs(&a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a))
	}

	if got := a["name"].Value.Type; got != openapi.TypeString {
		t.Fatalf("expected type %q, got %q", openapi.TypeString, got)
	}
}
