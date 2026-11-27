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

func TestParameter_Extensions(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name, a, b, want, err string
	}{
		{name: "none"},
		{name: "kept", a: `{"x-go-name":"QueryID"}`, want: `{"x-go-name":"QueryID"}`},
		{name: "taken", b: `{"x-go-name":"QueryID"}`, want: `{"x-go-name":"QueryID"}`},
		{name: "same", a: `{"x-go-name":"QueryID"}`, b: `{"x-go-name": "QueryID"}`, want: `{"x-go-name":"QueryID"}`},
		{
			name: "added",
			a:    `{"x-go-name":"QueryID","x-a":{"b":1,"c":2}}`, b: `{"x-z":true,"x-a":{"c":2,"b":1}}`,
			want: `{"x-go-name":"QueryID","x-a":{"b":1,"c":2},"x-z":true}`,
		},
		{
			name: "conflict",
			a:    `{"x-go-name":"QueryID"}`, b: `{"x-go-name":"ReferenceCode"}`,
			err: `x-go-name: "QueryID" vs. "ReferenceCode"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := &openapi.Parameter{
				Name: "q", In: openapi.ParameterLocationQuery,
				Schema: &openapi.Schema{Type: openapi.TypeInteger}, Extensions: openapi.Extensions(tc.a),
			}
			b := &openapi.Parameter{
				Name: "q", In: openapi.ParameterLocationQuery,
				Schema: &openapi.Schema{Type: openapi.TypeInteger}, Extensions: openapi.Extensions(tc.b),
			}

			err := merge.Parameter(a, b)
			switch {
			case tc.err != "":
				if err == nil || err.Error() != tc.err {
					t.Fatalf("got error %v, want %s", err, tc.err)
				}

				return
			case err != nil:
				t.Fatal(err)
			}

			if string(a.Extensions) != tc.want {
				t.Fatalf("got extensions %s, want %s", a.Extensions, tc.want)
			}
		})
	}
}
