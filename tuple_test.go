package merge_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

// schemaOf decodes a schema from JSON.
func schemaOf(t *testing.T, data string) *openapi.Schema {
	t.Helper()

	s := &openapi.Schema{}
	if err := json.Unmarshal([]byte(data), s); err != nil {
		t.Fatal(err)
	}

	return s
}

func TestSchema_Tuples(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		a       string
		samples []string
		want    string
	}{
		{
			// pixellab's 3D keypoint: a recorded [x, y, z] is inferred as a list of numbers, as it keeps no length
			name:    "a list fitting every position is the tuple",
			a:       `{"type":"array","prefixItems":[{"type":"number","x-go-name":"X"},{"type":"number","x-go-name":"Y"},{"type":"number","x-go-name":"Z"}]}`,
			samples: []string{`{"type":"array","items":{"type":"integer"}}`, `{"type":"array","items":{"type":"number"}}`, `{"type":"array","items":{"type":"number"}}`},
			want:    `{"type":"array","prefixItems":[{"type":"number","format":"double","x-go-name":"X"},{"type":"number","format":"double","x-go-name":"Y"},{"type":"number","format":"double","x-go-name":"Z"}]}`,
		},
		{
			// before, each list was merged into the tuple, the first array, nesting a union around it every time
			name:    "a list goes to the list",
			a:       `{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"}]}`,
			samples: []string{`{"type":"array","items":{"type":"boolean"}}`, `{"type":"array","items":{"type":"boolean"}}`, `{"type":"array","items":{"type":"boolean"}}`},
			want:    `{"oneOf":[{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"}]},{"type":"array","items":{"type":"boolean"}}]}`,
		},
		{
			name:    "a tuple of another length is an alternative of its own",
			a:       `{"oneOf":[{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"}]},{"type":"array","items":{"type":"boolean"}}]}`,
			samples: []string{`{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"},{"type":"boolean"}]}`, `{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"},{"type":"boolean"}]}`},
			want:    `{"oneOf":[{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"}]},{"type":"array","items":{"type":"boolean"}},{"type":"array","prefixItems":[{"type":"string"},{"type":"integer"},{"type":"boolean"}]}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := schemaOf(t, tc.a)
			for _, sample := range tc.samples {
				if err := merge.Schema(a, schemaOf(t, sample), false); err != nil {
					t.Fatal(err)
				}
			}

			got, err := json.Marshal(a)
			if err != nil {
				t.Fatal(err)
			}

			want, err := json.Marshal(schemaOf(t, tc.want))
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != string(want) {
				t.Errorf("got  %s\nwant %s", got, want)
			}
		})
	}
}
