```bash
go get github.com/MarkRosemaker/openapi-merge
```


```go
import (
    "github.com/MarkRosemaker/openapi"
    merge "github.com/MarkRosemaker/openapi-merge"
)

// b is merged into a; a is modified in place.
if err := merge.Schema(a, b, false); err != nil {
    log.Fatal(err) // e.g. properties["age"].type: "string" != "integer"
}
```

The final argument to `Schema` marks whether the schemas describe a *parameter*,
which enables the scalar-or-array reconciliation above — that mismatch is an
artifact of how query parameters get sampled, and applying it to a request body
would mask a genuine conflict.

Merging is available for each object kind:

| Function | Merges |
|---|---|
| `merge.Schema(a, b *openapi.Schema, isParam bool)` | Two schemas |
| `merge.SchemaRefs(a *openapi.SchemaRefs, b openapi.SchemaRefs)` | Two sets of named schemas |
| `merge.Parameter(a, b *openapi.Parameter)` | Two parameters |
| `merge.Response(a, b *openapi.Response)` | Two responses |
| `merge.MediaType(a, b *openapi.MediaType)` | Two media types |
| `merge.Content(a *openapi.Content, b openapi.Content)` | Two content maps |

Errors carry the full JSON path to the conflict, via
[`errpath`](https://github.com/MarkRosemaker/errpath).
