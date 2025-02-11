package merge

import (
	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/openapi"
)

func SchemaRefs(a, b openapi.SchemaRefs) error {
	for keyB, sB := range b.ByIndex() {
		sA, ok := a[keyB]
		if !ok {
			a.Set(keyB, sB) // add the property
			continue
		}

		// merge the properties
		if err := Schema(sA.Value, sB.Value); err != nil {
			return &errpath.ErrKey{Key: keyB, Err: err}
		}
	}

	return nil
}
