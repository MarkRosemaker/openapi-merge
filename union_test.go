package merge_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
)

// object builds an object schema with the given properties, in order.
func object(props ...any) *openapi.Schema {
	s := &openapi.Schema{Type: openapi.TypeObject, Properties: openapi.Schemas{}}
	for i := 0; i < len(props); i += 2 {
		s.Properties.Set(props[i].(string), props[i+1].(*openapi.Schema))
	}

	return s
}

func pinned(v string) *openapi.Schema {
	return &openapi.Schema{Type: openapi.TypeString, Const: jsontext.Value(`"` + v + `"`)}
}

// sample is a string property as it is inferred from a recorded value.
func sample(v string) *openapi.Schema {
	return &openapi.Schema{Type: openapi.TypeString, Example: jsontext.Value(`"` + v + `"`)}
}

func refTo(name string, v *openapi.Schema) *openapi.Schema {
	return &openapi.Schema{Ref: &openapi.SchemaRef{Identifier: "#/components/schemas/" + name, Value: v}}
}

// variants is a tagged union of objects, told apart by their type, the way Notion's are.
func variants() (number, sel *openapi.Schema, union openapi.SchemaList) {
	number = object("type", pinned("number"), "number", object())
	number.Description = "A number property."
	sel = object("type", pinned("select"), "select", object())

	return number, sel, openapi.SchemaList{refTo("Number", number), refTo("Select", sel)}
}

func TestSchema_OneOfPicksVariantByPinnedValue(t *testing.T) {
	t.Parallel()

	number, sel, union := variants()
	a := &openapi.Schema{OneOf: union}

	if err := merge.Schema(a, object("type", sample("select"), "select", object("options", sample("x"))), false); err != nil {
		t.Fatal(err)
	}

	if _, ok := sel.Properties["select"].Properties["options"]; !ok {
		t.Error("the select variant did not get the recorded options")
	}

	if len(number.Properties["number"].Properties) != 0 {
		t.Error("the number variant got what was recorded for select")
	}

	if number.Description != "A number property." {
		t.Errorf("the number variant's description is %q, want it kept", number.Description)
	}
}

func TestSchema_AllOfOfCommonAndUnion(t *testing.T) {
	t.Parallel()

	common := object("id", &openapi.Schema{Type: openapi.TypeString})
	_, sel, union := variants()
	a := &openapi.Schema{AllOf: openapi.SchemaList{refTo("Common", common), {OneOf: union}}}

	b := object("id", sample("abc"), "name", sample("Tags"), "type", sample("select"), "select", object())
	if err := merge.Schema(a, b, false); err != nil {
		t.Fatal(err)
	}

	if common.Properties["id"].Example == nil {
		t.Error("the common part did not get its recorded property")
	}

	if _, ok := sel.Properties["name"]; !ok {
		t.Error("the select variant did not get the property no part declares")
	}

	if _, ok := sel.Properties["id"]; ok {
		t.Error("the select variant got a property the common part declares")
	}
}

func TestSchema_AnyOfOfUnions(t *testing.T) {
	t.Parallel()

	page := object("object", pinned("page"), "url", &openapi.Schema{Type: openapi.TypeString})
	dataSource := object("object", pinned("data_source"), "title", &openapi.Schema{Type: openapi.TypeString})
	a := &openapi.Schema{AnyOf: openapi.SchemaList{
		{AnyOf: openapi.SchemaList{refTo("Page", page)}},
		// a reference to a reference stands for what the last one refers to
		{AnyOf: openapi.SchemaList{refTo("DataSourceAlias", refTo("DataSource", dataSource))}},
	}}

	if err := merge.Schema(a, object("object", sample("data_source"), "is_inline", &openapi.Schema{
		Type: openapi.TypeBoolean,
	}), false); err != nil {
		t.Fatal(err)
	}

	if _, ok := dataSource.Properties["is_inline"]; !ok {
		t.Error("the data source did not get the recorded property")
	}

	if _, ok := page.Properties["is_inline"]; ok {
		t.Error("the page got what was recorded for a data source")
	}
}

func TestSchema_AlternativeOfAWiderType(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		alt, b *openapi.Schema
	}{
		"an integer is a number": {
			&openapi.Schema{Type: openapi.TypeNumber},
			&openapi.Schema{Type: openapi.TypeInteger, Example: jsontext.Value(`3`)},
		},
		"a string without a format takes any": {
			&openapi.Schema{Type: openapi.TypeString},
			&openapi.Schema{Type: openapi.TypeString, Format: openapi.FormatURI},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			a := &openapi.Schema{OneOf: openapi.SchemaList{tc.alt, {Type: openapi.TypeNull}}}
			if err := merge.Schema(a, tc.b, false); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSchema_NoVariantMatches(t *testing.T) {
	t.Parallel()

	_, _, union := variants()

	err := merge.Schema(&openapi.Schema{OneOf: union}, object("type", sample("rollup")), false)
	if err == nil {
		t.Fatal("expected an error")
	}

	if want := `oneOf: no branch matches type: "rollup"`; err.Error() != want {
		t.Errorf("got error %q, want %q", err, want)
	}
}

func TestSchema_AnyOfOnEitherSide(t *testing.T) {
	t.Parallel()

	// a sample merged into a union routes the same whichever side the union is on
	_, sel, union := variants()
	b := &openapi.Schema{AnyOf: union}
	a := object("type", sample("select"), "select", object("options", sample("x")))

	if err := merge.Schema(a, b, false); err != nil {
		t.Fatal(err)
	}

	if len(a.AnyOf) != 2 {
		t.Fatalf("the result is not the union, got %d alternatives", len(a.AnyOf))
	}

	if _, ok := sel.Properties["select"].Properties["options"]; !ok {
		t.Error("the select variant did not get the recorded options")
	}

	// an anyOf merged into a union merges each of its alternatives
	_, sel2, union2 := variants()
	if err := merge.Schema(&openapi.Schema{OneOf: union2}, &openapi.Schema{AnyOf: openapi.SchemaList{
		object("type", sample("select"), "select", object("color", sample("red"))),
	}}, false); err != nil {
		t.Fatal(err)
	}

	if _, ok := sel2.Properties["select"].Properties["color"]; !ok {
		t.Error("the select variant did not get the anyOf's alternative")
	}
}

func TestSchema_NoExampleBesideConst(t *testing.T) {
	t.Parallel()

	_, sel, union := variants()

	if err := merge.Schema(&openapi.Schema{OneOf: union}, object("type", sample("select"), "select", object()), false); err != nil {
		t.Fatal(err)
	}

	if ex := sel.Properties["type"].Example; ex != nil {
		t.Errorf("the pinned type got the example %s, which only repeats its const", ex)
	}
}

func TestSchema_ReferenceCycle(t *testing.T) {
	t.Parallel()

	// a reference that leads back to itself ends rather than looping
	loop := &openapi.Schema{}
	loop.Ref = &openapi.SchemaRef{Identifier: "#/components/schemas/Loop", Value: loop}

	_ = merge.Schema(&openapi.Schema{OneOf: openapi.SchemaList{loop}}, object(), false)
}
