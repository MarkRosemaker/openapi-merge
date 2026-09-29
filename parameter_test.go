package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestParameter(t *testing.T) {
	t.Parallel()

	a := &openapi.Parameter{
		Name:   "TODO",
		In:     openapi.ParameterLocationQuery,
		Schema: &openapi.Schema{Type: openapi.TypeString},
	}
	b := &openapi.Parameter{
		Name:        "id",
		In:          openapi.ParameterLocationQuery,
		Description: "the id",
		Schema:      &openapi.Schema{Type: openapi.TypeString},
	}

	if err := merge.Parameter(a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// results must be merged into a, the first argument
	if a.Name != "id" {
		t.Fatalf("expected name %q, got %q", "id", a.Name)
	}

	if a.Description != "the id" {
		t.Fatalf("expected description %q, got %q", "the id", a.Description)
	}
}
