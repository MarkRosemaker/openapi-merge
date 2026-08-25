package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestResponse(t *testing.T) {
	t.Parallel()

	// a starts out with no content at all, b introduces a new media type.
	// This must end up merged into a, even though a.Content starts out nil.
	a := &openapi.Response{Description: "TODO"}
	b := &openapi.Response{
		Description: "the response",
		Content: openapi.Content{
			"application/json": &openapi.MediaType{
				Schema: &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
			},
		},
	}

	if err := merge.Response(a, b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.Description != "the response" {
		t.Fatalf("expected description %q, got %q", "the response", a.Description)
	}

	if len(a.Content) != 1 {
		t.Fatalf("expected 1 content entry, got %d", len(a.Content))
	}

	if got := a.Content["application/json"].Schema.Value.Type; got != openapi.TypeString {
		t.Fatalf("expected type %q, got %q", openapi.TypeString, got)
	}
}
