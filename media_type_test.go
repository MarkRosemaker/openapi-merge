package merge_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
)

func TestMediaType(t *testing.T) {
	t.Parallel()

	a := &openapi.MediaType{}
	b := &openapi.MediaType{
		Schema:  &openapi.SchemaRef{Value: &openapi.Schema{Type: openapi.TypeString}},
		Example: jsontext.Value(`"foo"`),
	}

	require.NoError(t, merge.MediaType(a, b))

	// the result must land in a, the first argument
	require.Equal(t, openapi.TypeString, a.Schema.Value.Type)
	require.Equal(t, jsontext.Value(`"foo"`), a.Example)
}
