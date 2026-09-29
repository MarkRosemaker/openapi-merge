package merge_test

import (
	"reflect"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestSchema_NullAndArrayBounds(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		a, b *openapi.Schema
		want *openapi.Schema
	}{
		{
			"null and null",
			&openapi.Schema{Type: openapi.TypeNull},
			&openapi.Schema{Type: openapi.TypeNull},
			&openapi.Schema{Type: openapi.TypeNull},
		},
		{
			"nullable is kept",
			&openapi.Schema{Type: openapi.TypeString, Nullable: true},
			&openapi.Schema{Type: openapi.TypeString},
			&openapi.Schema{Type: openapi.TypeString, Nullable: true},
		},
		{
			"nullable is adopted",
			&openapi.Schema{Type: openapi.TypeString},
			&openapi.Schema{Type: openapi.TypeString, Nullable: true},
			&openapi.Schema{Type: openapi.TypeString, Nullable: true},
		},
		{
			"only ever seen empty, then with items",
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
			&openapi.Schema{Type: openapi.TypeArray, Items: &openapi.Schema{Type: openapi.TypeInteger}},
			&openapi.Schema{Type: openapi.TypeArray, Items: &openapi.Schema{Type: openapi.TypeInteger}},
		},
		{
			"with items, then seen empty",
			&openapi.Schema{Type: openapi.TypeArray, Items: &openapi.Schema{Type: openapi.TypeInteger}},
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
			&openapi.Schema{Type: openapi.TypeArray, Items: &openapi.Schema{Type: openapi.TypeInteger}},
		},
		{
			"always empty",
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
		},
		{
			"bounds widen",
			&openapi.Schema{Type: openapi.TypeArray, MinItems: 2, MaxItems: new(uint(3)), Items: &openapi.Schema{Type: openapi.TypeString}},
			&openapi.Schema{Type: openapi.TypeArray, MinItems: 1, MaxItems: new(uint(5)), Items: &openapi.Schema{Type: openapi.TypeString}},
			&openapi.Schema{Type: openapi.TypeArray, MinItems: 1, MaxItems: new(uint(5)), Items: &openapi.Schema{Type: openapi.TypeString}},
		},
		{
			"null, then seen empty",
			&openapi.Schema{Type: openapi.TypeNull},
			&openapi.Schema{Type: openapi.TypeArray, MaxItems: new(uint(0))},
			&openapi.Schema{Type: openapi.TypeArray, Nullable: true, MaxItems: new(uint(0))},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := merge.Schema(tc.a, tc.b, false); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.a, tc.want) {
				t.Fatalf("got %+v, want %+v", tc.a, tc.want)
			}

			if err := tc.a.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
