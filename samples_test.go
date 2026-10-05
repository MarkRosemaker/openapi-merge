package merge_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestSchema_Samples(t *testing.T) {
	t.Parallel()

	paragraph := `{"type":"object","properties":{"type":{"type":"string","const":"paragraph"},"paragraph":{"type":"string"}}}`
	heading := `{"type":"object","properties":{"type":{"type":"string","const":"heading"},"heading":{"type":"string"}}}`
	paragraphSample := `{"type":"object","properties":{"type":{"type":"string","example":"paragraph"},"paragraph":{"type":"string"}}}`
	headingSample := `{"type":"object","properties":{"type":{"type":"string","example":"heading"},"heading":{"type":"string"}}}`

	for _, tc := range []struct {
		name    string
		a       string
		samples []string
		want    string
	}{
		{
			// merged into one, the two samples would hold both members and match neither alternative
			name:    "each sample into the alternative it matches",
			a:       `{"oneOf":[` + paragraph + `,` + heading + `]}`,
			samples: []string{paragraphSample, headingSample},
			want:    `{"oneOf":[` + paragraph + `,` + heading + `]}`,
		},
		{
			name:    "each sample in turn into a schema that is no union",
			a:       `{"type":"object","properties":{"type":{"type":"string"}}}`,
			samples: []string{paragraphSample, headingSample},
			want:    `{"type":"object","properties":{"type":{"type":"string","example":"paragraph"},"paragraph":{"type":"string"},"heading":{"type":"string"}}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			samples := make([]*openapi.Schema, len(tc.samples))
			for i, s := range tc.samples {
				samples[i] = schemaOf(t, s)
			}

			a := schemaOf(t, tc.a)
			if err := merge.Schema(a, merge.Samples(samples...), false); err != nil {
				t.Fatal(err)
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

func TestSchema_SamplesJoinSamples(t *testing.T) {
	t.Parallel()

	// samples are no alternatives to route among: more of them join the union, for the caller to collapse
	a := merge.Samples(schemaOf(t, `{"type":"string"}`))
	if err := merge.Schema(a, merge.Samples(schemaOf(t, `{"type":"integer"}`)), false); err != nil {
		t.Fatal(err)
	}

	if err := merge.Schema(a, schemaOf(t, `{"type":"boolean"}`), false); err != nil {
		t.Fatal(err)
	}

	if !merge.IsSamples(a) || len(a.AnyOf) != 3 {
		t.Errorf("got %+v, want three samples", a)
	}

	if merge.IsSamples(schemaOf(t, `{"anyOf":[{"type":"string"}]}`)) {
		t.Error("a plain anyOf taken for samples")
	}
}
