package merge_test

import (
	"encoding/json/v2"
	"testing"

	merge "github.com/MarkRosemaker/openapi-merge"
)

func TestSchema_DateOrDateTime(t *testing.T) {
	t.Parallel()

	const want = `{"oneOf":[{"type":"string","format":"date","example":"2026-10-05"},{"type":"string","format":"date-time","example":"2026-10-05T10:00:00Z"}]}`

	// either way round, and further samples of either go to their own
	for _, order := range [][2]string{
		{`{"type":"string","format":"date","example":"2026-10-05"}`, `{"type":"string","format":"date-time","example":"2026-10-05T10:00:00Z"}`},
		{`{"type":"string","format":"date-time","example":"2026-10-05T10:00:00Z"}`, `{"type":"string","format":"date","example":"2026-10-05"}`},
	} {
		a := schemaOf(t, order[0])
		if err := merge.Schema(a, schemaOf(t, order[1]), false); err != nil {
			t.Fatal(err)
		}

		for _, more := range order {
			if err := merge.Schema(a, schemaOf(t, more), false); err != nil {
				t.Fatal(err)
			}
		}

		got, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}

		if string(got) != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	}
}
