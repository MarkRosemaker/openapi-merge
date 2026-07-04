package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
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

	require.NoError(t, merge.Content(&a, b))
	require.Len(t, a, 1)
	require.Same(t, mt, a["application/json"])
}

func TestContent_MergeExistingKey(t *testing.T) {
	t.Parallel()

	a := openapi.Content{"application/json": &openapi.MediaType{
		Schema: &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
	}}
	b := openapi.Content{"application/json": &openapi.MediaType{
		Schema: &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
	}}

	require.NoError(t, merge.Content(&a, b))
	require.Len(t, a, 1)
	require.Equal(t, openapi.TypeString, a["application/json"].Schema.Value.Type)
}
