package merge_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

// TestSchema_Golden merges, into the property a of each case in testdata/before.json, every schema of its prefixItems
// in turn, and must fail with its x-error where it has one. The merged schemas are then dropped, and a failed case
// with them, and the document must be testdata/after.json, which shows each result and what it changed in the
// schemas a case refers to. Both files are edited by hand.
func TestSchema_Golden(t *testing.T) {
	t.Parallel()

	doc, err := openapi.LoadFromFile("testdata/before.json")
	if err != nil {
		t.Fatal(err)
	}

	for name, c := range doc.Components.Schemas.ByIndex() {
		a := c.Properties["a"]
		if a == nil {
			continue // a schema the cases refer to
		}

		var want struct {
			Error string `json:"x-error"`
		}

		if c.Extensions != nil {
			if err := json.Unmarshal(c.Extensions, &want); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
		}

		for _, b := range c.PrefixItems {
			if err = merge.Schema(a, b, false); err != nil {
				break
			}
		}

		switch {
		case want.Error != "":
			if err == nil || err.Error() != want.Error {
				t.Errorf("%s: got error %v, want %s", name, err, want.Error)
			}

			delete(doc.Components.Schemas, name)
		case err != nil:
			t.Errorf("%s: %v", name, err)
		}

		c.PrefixItems = nil
	}

	b, err := doc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	got := jsontext.Value(b)
	if err := got.Indent(jsontext.WithIndent("  ")); err != nil {
		t.Fatal(err)
	}

	got = append(got, '\n')

	want, err := os.ReadFile("testdata/after.json")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		gotLines, wantLines := bytes.Split(got, []byte("\n")), bytes.Split(want, []byte("\n"))
		for i := range min(len(gotLines), len(wantLines)) {
			if !bytes.Equal(gotLines[i], wantLines[i]) {
				t.Fatalf("after.json line %d: got %s, want %s", i+1, bytes.TrimSpace(gotLines[i]), bytes.TrimSpace(wantLines[i]))
			}
		}

		t.Fatalf("got %d lines, want %d", len(gotLines), len(wantLines))
	}
}
