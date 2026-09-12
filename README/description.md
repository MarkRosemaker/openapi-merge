---
tagline: Two views of the same endpoint, reconciled.
logo:
    alt: A gopher clicking two incomplete jigsaw pieces together into one complete piece
    source: openapi-merge.jpg
    width: 500
---

<div align="center" id=badges>

![Code Coverage](https://img.shields.io/badge/coverage-65.7%25-yellowgreen)

</div>





`openapi-merge` combines two [OpenAPI 3.x](https://spec.openapis.org/oas/v3.1.0)
objects into one that covers both. It exists for the problem of *incomplete
evidence*: when a schema is inferred from a sample of real data, each sample tells
you only part of the story, and the parts have to be reconciled.
