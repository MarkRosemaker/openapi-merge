package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
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

	require.NoError(t, merge.Response(a, b))

	require.Equal(t, "the response", a.Description)
	require.Len(t, a.Content, 1)
	require.Equal(t, openapi.TypeString, a.Content["application/json"].Schema.Value.Type)
}
