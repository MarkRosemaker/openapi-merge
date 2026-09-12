Observe `GET /users/{id}` once and you might see `{"id": 1, "name": "Alice"}`.
Observe it again and you get `{"id": 2, "name": "Bob", "nickname": null}`. Neither
response is the schema. The schema is what you get by merging them: three
properties, one of them optional, one of them of unknown type.

That is what this module does. It is used by
[`openapi-enrich`](https://github.com/MarkRosemaker/openapi-enrich), which builds
specifications from recorded HTTP traffic and calls in here every time a second
observation of the same endpoint arrives.

Merging is destructive and asymmetric by design: `b` is merged **into** `a`, in
place. If the two cannot be reconciled, an error is returned describing the exact
JSON path at which they conflict.
