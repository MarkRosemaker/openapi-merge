package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
)

func TestSchemaRefs(t *testing.T) {
	t.Parallel()

	// a starts out nil: merging must still make the new entry visible to the
	// caller, even though maps are normally passed by value.
	var a openapi.SchemaRefs

	sr := &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}}
	b := openapi.SchemaRefs{"name": sr}

	require.NoError(t, merge.SchemaRefs(&a, b))
	require.Len(t, a, 1)
	require.Same(t, sr, a["name"])
}

func TestSchemaRefs_MergeExistingKey(t *testing.T) {
	t.Parallel()

	a := openapi.SchemaRefs{"name": &openapi.SchemaRef{
		Value: &openapi.Schema{Type: openapi.TypeString},
	}}
	b := openapi.SchemaRefs{"name": &openapi.SchemaRef{
		Value: &openapi.Schema{Type: openapi.TypeString},
	}}

	require.NoError(t, merge.SchemaRefs(&a, b))
	require.Len(t, a, 1)
	require.Equal(t, openapi.TypeString, a["name"].Value.Type)
}
