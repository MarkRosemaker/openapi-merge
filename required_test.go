package merge_test

import (
	"encoding/json/v2"
	"testing"

	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestSchema_Required(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name, a, b, want string
	}{
		{
			// a value without a member makes it optional; one b requires but a does not stays optional
			name: "only what both require",
			a:    `{"type":"object","required":["id","name"],"properties":{"id":{"type":"string"},"name":{"type":"string"}}}`,
			b:    `{"type":"object","required":["id","nick"],"properties":{"id":{"type":"string"},"nick":{"type":"string"}}}`,
			want: `{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"nick":{"type":"string"}},"required":["id"]}`,
		},
		{
			// Notion's filter forms pin nothing; the one declaring the sample's member takes it, not the first
			name: "the branch declaring the member",
			a: `{"anyOf":[` +
				`{"type":"object","required":["or"],"properties":{"or":{"type":"array","items":{"type":"string"}}}},` +
				`{"type":"object","required":["and"],"properties":{"and":{"type":"array","items":{"type":"string"}}}}]}`,
			b: `{"type":"object","required":["and"],"properties":{"and":{"type":"array","items":{"type":"string"}}}}`,
			want: `{"anyOf":[` +
				`{"type":"object","properties":{"or":{"type":"array","items":{"type":"string"}}},"required":["or"]},` +
				`{"type":"object","properties":{"and":{"type":"array","items":{"type":"string"}}},"required":["and"]}]}`,
		},
		{
			// Notion's database lacks database_type, which the full form requires: it is still the full form
			name: "the fuller branch, made optional",
			a: `{"oneOf":[` +
				`{"type":"object","required":["object","id"],"properties":{"object":{"type":"string","const":"database"},"id":{"type":"string"}}},` +
				`{"type":"object","required":["object","id","title","database_type"],"properties":{"object":{"type":"string","const":"database"},"id":{"type":"string"},"title":{"type":"string"},"database_type":{"type":"string"}}}]}`,
			b: `{"type":"object","required":["object","id","title"],"properties":{"object":{"type":"string","example":"database"},"id":{"type":"string"},"title":{"type":"string"}}}`,
			want: `{"oneOf":[` +
				`{"type":"object","properties":{"object":{"type":"string","const":"database"},"id":{"type":"string"}},"required":["object","id"]},` +
				`{"type":"object","properties":{"object":{"type":"string","const":"database"},"id":{"type":"string"},"title":{"type":"string"},"database_type":{"type":"string"}},"required":["object","id","title"]}]}`,
		},
		{
			// split among an allOf's parts, a sample narrows each part only by what it leaves out
			name: "parts of an allOf",
			a: `{"allOf":[` +
				`{"type":"object","required":["plain_text"],"properties":{"plain_text":{"type":"string"}}},` +
				`{"oneOf":[{"type":"object","required":["type","text"],"properties":{"type":{"type":"string","const":"text"},"text":{"type":"string"}}}]}]}`,
			b: `{"type":"object","required":["plain_text","type","text"],"properties":{"plain_text":{"type":"string"},"type":{"type":"string","example":"text"},"text":{"type":"string"}}}`,
			want: `{"allOf":[` +
				`{"type":"object","properties":{"plain_text":{"type":"string"}},"required":["plain_text"]},` +
				`{"oneOf":[{"type":"object","properties":{"type":{"type":"string","const":"text"},"text":{"type":"string"}},"required":["type","text"]}]}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := schemaOf(t, tc.a)
			if err := merge.Schema(a, schemaOf(t, tc.b), false); err != nil {
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
