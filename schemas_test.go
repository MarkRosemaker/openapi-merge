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
	var a openapi.Schemas

	sr := &openapi.Schema{Type: openapi.TypeString}
	b := openapi.Schemas{"name": sr}

	if err := merge.Schemas(&a, b); err != nil {
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

	a := openapi.Schemas{"name": &openapi.Schema{Type: openapi.TypeString}}
	b := openapi.Schemas{"name": &openapi.Schema{Type: openapi.TypeString}}

	if err := merge.Schemas(&a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a))
	}

	if got := a["name"].Type; got != openapi.TypeString {
		t.Fatalf("expected type %q, got %q", openapi.TypeString, got)
	}
}
