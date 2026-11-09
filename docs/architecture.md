# Architecture

`Schema` merges two schemas of one value into one that describes both: the
specification so far, such as an upstream one or what earlier samples made of
it, and the schema of a newly recorded sample, as openapi-enrich derives it.
The other entry points (`Content`, `MediaType`, `Response`, `Parameter`,
`Schemas`) walk down to it.

## Merging only widens

A merged schema may describe more values than either side, never fewer. A
schema that rejects a value the API sends breaks every client generated from
it, so where the sides disagree, the merge keeps what both allow: an optional
member, a nullable type, a wider integer, a number for an integer, a `oneOf`
of shapes that do not merge.

Narrowing is openapi-enrich's: it derives a schema from recorded values alone,
and where every value it has seen shares a format, that format is the best it
knows, until a value that does not fit widens it again here.

### Formats

`reconcileFormats` decides a format. A string with a format merged with one
without keeps none (`default` in its string case). The side without one is a
contract that allows any string, such as an upstream specification declares,
or what values that share no format left behind, and nothing tells the two
apart. Taking the other side's format, such as the `uuid` of one recorded ID,
would claim of every value what one value showed: go-api-libs/notion's error
details, declared as strings, were recorded once as a UUID, and a client that
required one would fail on `"none"`. The golden case
`bare-string-keeps-no-format` holds this.

The exceptions are where both sides' values fit the format: `date-basic`
stays when the value without it is such a date too, and `date` or
`date-time` stays when that value is an empty string, a date left out.

An integer without a format is the exception that widens: every sample is
recorded without one, and an integer of no stated width already allows any,
so the declared width stays, the wider of two.
