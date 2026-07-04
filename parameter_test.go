package merge_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
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

	require.NoError(t, merge.Parameter(a, b))

	// results must be merged into a, the first argument
	require.Equal(t, "id", a.Name)
	require.Equal(t, "the id", a.Description)
}
