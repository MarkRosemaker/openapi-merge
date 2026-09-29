Beyond combining properties and widening optionality, the merge handles the
particular ways that sample-derived schemas disagree:

- **Null** — a value observed only as `null` has the type `null`. Merged with a
  real type, the result is that type, made nullable (`["string", "null"]`), rather
  than a conflict.
- **Arrays only ever seen empty** — `{"type": "array", "maxItems": 0}` says nothing
  about the items, so the other side's items are adopted, and item bounds widen to
  cover both sides.
- **Numeric widening** — an integer in one sample and a floating-point number in
  another merge to a number.
- **Dates in two encodings** — a value seen as a date-time string in one sample and
  as a Unix timestamp integer in another becomes a `oneOf` of the two, rather than
  one silently discarding the other.
- **`oneOf` routing** — when one side already covers several shapes, the other is
  merged into whichever branch it matches.
- **Scalar-or-array parameters** — a parameter that appeared as a bare value in one
  sample and as an array in another merges the value into the array's item schema.
- **Enums** — an example value not yet present in an enum is added to it.
