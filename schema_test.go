package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

// TestSchema_ReferenceCycle: a reference that leads back to itself ends rather than looping. A loaded document cannot
// hold one, so it is built here.
func TestSchema_ReferenceCycle(t *testing.T) {
	t.Parallel()

	loop := &openapi.Schema{}
	loop.Ref = &openapi.SchemaRef{Identifier: "#/components/schemas/Loop", Value: loop}

	_ = merge.Schema(&openapi.Schema{OneOf: openapi.SchemaList{loop}}, &openapi.Schema{Type: openapi.TypeObject}, false)
}

// TestIsSamples: only a union marked as samples is one.
func TestIsSamples(t *testing.T) {
	t.Parallel()

	if !merge.IsSamples(merge.Samples(&openapi.Schema{Type: openapi.TypeString})) {
		t.Error("Samples are not samples")
	}

	if merge.IsSamples(&openapi.Schema{AnyOf: openapi.SchemaList{{Type: openapi.TypeString}}}) {
		t.Error("a plain anyOf taken for samples")
	}
}
