Beyond combining properties and widening optionality, the merge handles the
particular ways that sample-derived schemas disagree:

- **Required** — an object requires only what both sides require: a value without
  a member makes it optional, as a recorded sample requires every member it has.
  What an inline `allOf` part requires is narrowed the same way.
- **Null** — a value observed only as `null` has the type `null`. Merged with a
  real type, the result is that type, made nullable (`["string", "null"]`), rather
  than a conflict.
- **Arrays only ever seen empty** — `{"type": "array", "maxItems": 0}` says nothing
  about the items, so the other side's items are adopted, and item bounds widen to
  cover both sides.
- **Tuples** — a tuple (`prefixItems`) merges position by position with a tuple of
  its length. A sample keeps no length, so a list whose item type fits every
  position, such as `[x, y, z]` recorded as a list of numbers, merges into each
  position too. Any other list or tuple becomes an alternative beside it in a
  `oneOf`, and later samples go to the alternative of their own shape, or add one.
- **Numeric widening** — an integer in one sample and a floating-point number in
  another merge to a number.
- **Dates in two encodings** — a value seen as a date-time string in one sample and
  as a Unix timestamp integer in another becomes a `oneOf` of the two, rather than
  one silently discarding the other.
- **Dates with and without a time** — a value seen as a `date` in one sample and
  as a `date-time` in another, as Notion's date `start` is, becomes a `oneOf` of
  the two.
- **Union routing** — when one side already covers several shapes, with `oneOf`
  or `anyOf`, the other is merged into whichever branch it matches. Of a tagged
  union's objects, the branch is the one whose pinned properties — a `const` or
  one-value `enum`, such as `"type": {"const": "select"}` — the sample has, so a
  sample never lands in a sibling's branch; of several that match, it is the one
  that declares the most of the sample's properties, so a full object wins over
  its partial form. Of branches that pin nothing, it is likewise the one that
  declares the most of the sample's properties, and, of those that declare as
  many, one whose required properties the sample has. With none that matches, the merge
  fails and names the values it got. A branch that is itself a union matches as
  its own branches do, an integer matches a number branch, and a string without
  a format matches when no branch has the sample's.
- **Samples** — a union marked `x-samples` (see `Samples`) holds samples of one
  value, such as the elements of a recorded array, not alternatives the value may
  take. Merged into a union, each sample goes into the branch it matches, so the
  elements of a list of mixed variants each reach their own; merged into a schema
  that is no union, each is merged into it in turn; merged into another union of
  samples, they join it, for the caller to collapse. Unlike `examples`, which
  hold values and stay in the specification, it holds schemas inferred from
  values and is only a working marker between the inference and the merge:
  openapi-enrich routes or collapses every one, so none reaches a finished
  specification.
- **Common properties beside a union** — an `allOf` of objects and one union
  takes each sampled property into the part that declares it, and the rest into
  the branch the whole sample matches.
- **Scalar-or-array parameters** — a parameter that appeared as a bare value in one
  sample and as an array in another merges the value into the array's item schema.
- **Enums** — an example value not yet present in an enum is added to it.
