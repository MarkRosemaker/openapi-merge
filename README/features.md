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
- **Union routing** — when one side already covers several shapes, with `oneOf`
  or `anyOf`, the other is merged into whichever branch it matches. Of a tagged
  union's objects, the branch is the one whose pinned properties — a `const` or
  one-value `enum`, such as `"type": {"const": "select"}` — the sample has, so a
  sample never lands in a sibling's branch; with none that matches, the merge
  fails and names the values it got. A branch that is itself a union matches as
  its own branches do, an integer matches a number branch, and a string without
  a format matches when no branch has the sample's.
- **Common properties beside a union** — an `allOf` of objects and one union
  takes each sampled property into the part that declares it, and the rest into
  the branch the whole sample matches.
- **Scalar-or-array parameters** — a parameter that appeared as a bare value in one
  sample and as an array in another merges the value into the array's item schema.
- **Enums** — an example value not yet present in an enum is added to it.
