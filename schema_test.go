package merge_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
	merge "github.com/MarkRosemaker/openapi-merge"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		a, b *openapi.Schema
		want *openapi.Schema
	}{
		{&openapi.Schema{
			Type: openapi.TypeString,
			Enum: []string{},
		}, &openapi.Schema{
			Type:    openapi.TypeString,
			Example: jsontext.Value(`"foo"`),
		}, &openapi.Schema{
			Type:    openapi.TypeString,
			Enum:    []string{"foo"},
			Example: jsontext.Value(`"foo"`),
		}},
		{&openapi.Schema{
			Type:    openapi.TypeString,
			Format:  openapi.FormatURI,
			Example: jsontext.Value(`"https://www.example.com/"`),
		}, &openapi.Schema{
			Type:    openapi.TypeObject,
			Example: jsontext.Value(`null`),
			// TODO: Nullable
		}, &openapi.Schema{
			Type:    openapi.TypeString,
			Format:  openapi.FormatURI,
			Example: jsontext.Value(`"https://www.example.com/"`),
			// TODO: Nullable
		}},
		{&openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.SchemaRefs{
				"type": &openapi.SchemaRef{Value: &openapi.Schema{
					Type:    openapi.TypeString,
					Example: jsontext.Value(`"text"`),
					Enum:    []string{"text"},
				}},
				"text": &openapi.SchemaRef{Value: &openapi.Schema{
					Type: openapi.TypeObject,
					Properties: openapi.SchemaRefs{
						"content": &openapi.SchemaRef{Value: &openapi.Schema{
							Type:    openapi.TypeString,
							Example: jsontext.Value(`"This is a simple paragraph."`),
						}},
						"link": &openapi.SchemaRef{Value: &openapi.Schema{
							Type: openapi.TypeObject,
							Properties: openapi.SchemaRefs{
								"url": &openapi.SchemaRef{Value: &openapi.Schema{
									Type:    openapi.TypeString,
									Format:  openapi.FormatURI,
									Example: jsontext.Value(`"https://www.example.com/"`),
								}},
							},
						}},
					},
					Required: []string{"content", "link"},
					Example:  jsontext.Value(`{"content":"This is a simple paragraph.","link":null}`),
				}},

				// a: {
				// "properties":{"annotations":{"type":"object","properties":{"bold":{"type":"boolean","example":false},"italic":{"type":"boolean","example":false},"strikethrough":{"type":"boolean","example":false},"underline":{"type":"boolean","example":false},"code":{"type":"boolean","example":false},"color":{"type":"string","example":"default"}},"required":["bold","italic","strikethrough","underline","code","color"],"example":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"}},"plain_text":{"type":"string","example":"This is a simple paragraph."},"href":{"type":"string","format":"uri","example":"https://www.example.com/"},"mention":{"type":"object","properties":{"type":{"type":"string","example":"link_mention"},"link_mention":{"type":"object","properties":{"href":{"type":"string","format":"uri","example":"https://example.com/"},"title":{"type":"string","example":"Example Domain"},"description":{"type":"string","example":"This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission."}},"required":["href","title","description"],"example":{"href":"https://example.com/","title":"Example Domain","description":"This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission."}},"database":{"type":"object","properties":{"id":{"type":"string","format":"uuid","example":"7a3c647e-4c1e-4c27-bf1d-cfb0105e55ce"}},"required":["id"],"example":{"id":"7a3c647e-4c1e-4c27-bf1d-cfb0105e55ce"}}},"required":["type","link_mention"],"example":{"type":"link_mention","link_mention":{"href":"https://example.com/","title":"Example Domain","description":"This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission."}}},"equation":{"type":"object","properties":{"expression":{"type":"string","example":"e^{\\pi i}+1=0"}},"required":["expression"],"example":{"expression":"e^{\\pi i}+1=0"}}}
			},

			Example: jsontext.Value(`{"type":"text","text":{"content":"This is a simple paragraph.","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"This is a simple paragraph.","href":null}`),
		}, &openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.SchemaRefs{
				"type": &openapi.SchemaRef{Value: &openapi.Schema{
					Type:    openapi.TypeString,
					Example: jsontext.Value(`"mention"`),
				}},
			},
			// b: {
			// "properties":{"mention":{"type":"object","properties":{"type":{"type":"string","example":"user"},"user":{"type":"object","properties":{"object":{"type":"string","example":"user"},"id":{"type":"string","format":"uuid","example":"af171d5d-c36f-45bc-a0a3-6086c0dafa45"},"name":{"type":"string","example":"Fae Tools"},"avatar_url":{"type":"string","format":"uri","example":"https://lh3.googleusercontent.com/a-/AOh14Gi54BUKkLrZ2IX8ORURI__6avK9zjCYXdhbmthj=s100"},"type":{"type":"string","example":"person"},"person":{"type":"object","properties":{"email":{"type":"string","format":"email","example":"mark@faetools.com"}},"required":["email"],"example":{"email":"mark@faetools.com"}}},"required":["object","id","name","avatar_url","type","person"],"example":{"object":"user","id":"af171d5d-c36f-45bc-a0a3-6086c0dafa45","name":"Fae Tools","avatar_url":"https://lh3.googleusercontent.com/a-/AOh14Gi54BUKkLrZ2IX8ORURI__6avK9zjCYXdhbmthj=s100","type":"person","person":{"email":"mark@faetools.com"}}}},"required":["type","user"],"example":{"type":"user","user":{"object":"user","id":"af171d5d-c36f-45bc-a0a3-6086c0dafa45","name":"Fae Tools","avatar_url":"https://lh3.googleusercontent.com/a-/AOh14Gi54BUKkLrZ2IX8ORURI__6avK9zjCYXdhbmthj=s100","type":"person","person":{"email":"mark@faetools.com"}}}},"annotations":{"type":"object","properties":{"bold":{"type":"boolean","example":false},"italic":{"type":"boolean","example":false},"strikethrough":{"type":"boolean","example":false},"underline":{"type":"boolean","example":false},"code":{"type":"boolean","example":false},"color":{"type":"string","example":"default"}},"required":["bold","italic","strikethrough","underline","code","color"],"example":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"}},"plain_text":{"type":"string","example":"@Fae Tools"},"href":{"type":"object"},"text":{"type":"object","properties":{"content":{"type":"string","example":" created this page."},"link":{"type":"object"}},"required":["content","link"],"example":{"content":" created this page.","link":null}}},"required":["type","mention","annotations","plain_text","href"],"example":{"type":"mention","mention":{"type":"user","user":{"object":"user","id":"af171d5d-c36f-45bc-a0a3-6086c0dafa45","name":"Fae Tools","avatar_url":"https://lh3.googleusercontent.com/a-/AOh14Gi54BUKkLrZ2IX8ORURI__6avK9zjCYXdhbmthj=s100","type":"person","person":{"email":"mark@faetools.com"}}},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"@Fae Tools","href":null}}

			Required: []string{"type", "mention", "annotations", "plain_text", "href"},
		}, &openapi.Schema{
			Type: openapi.TypeObject,
			Properties: openapi.SchemaRefs{
				"type": &openapi.SchemaRef{Value: &openapi.Schema{
					Type:    openapi.TypeString,
					Example: jsontext.Value(`"text"`),
					Enum:    []string{"text", "mention"},
				}},
				"text": &openapi.SchemaRef{Value: &openapi.Schema{
					Type: openapi.TypeObject,
					Properties: openapi.SchemaRefs{
						"content": &openapi.SchemaRef{Value: &openapi.Schema{
							Type:    openapi.TypeString,
							Example: jsontext.Value(`"This is a simple paragraph."`),
						}},
						"link": &openapi.SchemaRef{Value: &openapi.Schema{
							Type: openapi.TypeObject,
							Properties: openapi.SchemaRefs{
								"url": &openapi.SchemaRef{Value: &openapi.Schema{
									Type:    openapi.TypeString,
									Format:  openapi.FormatURI,
									Example: jsontext.Value(`"https://www.example.com/"`),
								}},
							},
						}},
					},
					Required: []string{"content", "link"},
					Example:  jsontext.Value(`{"content":"This is a simple paragraph.","link":null}`),
				}},
			},
			Example: jsontext.Value(`{"type":"text","text":{"content":"This is a simple paragraph.","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"This is a simple paragraph.","href":null}`),
		}},
		{
			&openapi.Schema{
				Description: "The archived status of the page.",
				Type:        openapi.TypeBoolean,
				Example:     jsontext.Value(`false`),
			},
			&openapi.Schema{
				Type:    openapi.TypeBoolean,
				Example: jsontext.Value(`true`),
			},
			&openapi.Schema{
				Description: "The archived status of the page.",
				Type:        openapi.TypeBoolean,
				Example:     jsontext.Value(`false`),
			},
		},
		{
			&openapi.Schema{
				Type: openapi.TypeArray,
				Items: &openapi.SchemaRef{
					Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props(
							"type", &openapi.Schema{
								Type:    openapi.TypeString,
								Example: jsontext.Value(`"foo"`),
								Enum:    []string{"foo"},
							},
							"foo", &openapi.Schema{
								Type:    openapi.TypeBoolean,
								Example: jsontext.Value("true"),
							}),
					},
				},
			},
			&openapi.Schema{
				Type: openapi.TypeArray,
				Items: &openapi.SchemaRef{
					Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props(
							"type", &openapi.Schema{
								Type:    openapi.TypeString,
								Example: jsontext.Value(`"bar"`),
							},
							"bar", &openapi.Schema{
								Type:    openapi.TypeBoolean,
								Example: jsontext.Value("true"),
							}),
					},
				},
			},
			&openapi.Schema{
				Type: openapi.TypeArray,
				Items: &openapi.SchemaRef{
					Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props(
							"type", &openapi.Schema{
								Type:    openapi.TypeString,
								Example: jsontext.Value(`"foo"`),
								Enum:    []string{"foo", "bar"},
							},
							"foo", &openapi.Schema{
								Type:    openapi.TypeBoolean,
								Example: jsontext.Value("true"),
							},
							"bar", &openapi.Schema{
								Type:    openapi.TypeBoolean,
								Example: jsontext.Value("true"),
							}),
					},
				},
			},
		},
		{
			&openapi.Schema{Type: openapi.TypeInteger, Example: jsontext.Value("100")},
			&openapi.Schema{Type: openapi.TypeInteger, Example: jsontext.Value("100")},
			&openapi.Schema{Type: openapi.TypeInteger, Example: jsontext.Value("100")},
		},
		{
			&openapi.Schema{
				Description: "a URL",
				Type:        openapi.TypeString,
				Format:      openapi.FormatURI,
				Example:     jsontext.Value(`"https://www.example.com/"`),
			},
			&openapi.Schema{
				Type:       openapi.TypeObject,
				Properties: openapi.SchemaRefs{},
				Example:    jsontext.Value("null"),
			},
			&openapi.Schema{
				Description: "a URL",
				Type:        openapi.TypeString,
				Format:      openapi.FormatURI,
				Example:     jsontext.Value(`"https://www.example.com/"`),
			},
		},
		{
			&openapi.Schema{
				Type:       openapi.TypeObject,
				Properties: openapi.SchemaRefs{},
				Example:    jsontext.Value("null"),
			},
			&openapi.Schema{
				Type:    openapi.TypeString,
				Format:  openapi.FormatURI,
				Example: jsontext.Value(`"https://www.example.com/"`),
			},
			&openapi.Schema{
				Type:    openapi.TypeString,
				Format:  openapi.FormatURI,
				Example: jsontext.Value(`"https://www.example.com/"`),
			},
		},
		{
			&openapi.Schema{
				Type:    openapi.TypeString,
				Example: jsontext.Value(`"some text"`),
			},
			&openapi.Schema{
				Type:    openapi.TypeString,
				Format:  openapi.FormatURI,
				Example: jsontext.Value(`"https://www.example.com/"`),
			},
			&openapi.Schema{
				Type:    openapi.TypeString,
				Example: jsontext.Value(`"some text"`),
			},
		},
		{
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
				},
			},
			&openapi.Schema{
				Type: openapi.TypeObject,
				Properties: props("bar", &openapi.Schema{
					Type: openapi.TypeInteger,
				}),
			},
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
				},
			},
		},
		{
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
				},
			},
			&openapi.Schema{
				Type: openapi.TypeObject,
				Properties: props("foo", &openapi.Schema{
					Type: openapi.TypeString,
				}),
			},
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
				},
			},
		},
		{
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
				},
			},
			&openapi.Schema{
				Type: openapi.TypeObject,
				Properties: props("baz", &openapi.Schema{
					Type: openapi.TypeBoolean,
				}),
			},
			&openapi.Schema{
				AllOf: openapi.SchemaRefList{
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("foo", &openapi.Schema{
							Type: openapi.TypeString,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("bar", &openapi.Schema{
							Type: openapi.TypeInteger,
						}),
					}},
					&openapi.SchemaRef{Value: &openapi.Schema{
						Type: openapi.TypeObject,
						Properties: props("baz", &openapi.Schema{
							Type: openapi.TypeBoolean,
						}),
					}},
				},
			},
		},

		{
			&openapi.Schema{
				Type:   openapi.TypeInteger,
				Format: openapi.FormatDate,
			},
			&openapi.Schema{
				Type: openapi.TypeInteger,
			},
			&openapi.Schema{
				Type:   openapi.TypeInteger,
				Format: openapi.FormatDate,
			},
		},
	} {
		require.NoError(t, merge.Schema(tc.a, tc.b))
		require.Equal(t, tc.want, tc.a)
		require.NoError(t, tc.a.Validate())
	}
}

func TestSchema_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		a, b *openapi.Schema
		err  string
	}{
		// {nil, nil, "schema a is nil"},
		// {&openapi.Schema{}, nil, "schema b is nil"},
	} {
		err := merge.Schema(tc.a, tc.b)
		require.Error(t, err)
		require.Equal(t, tc.err, err.Error())
	}
}

func props(keyVals ...any) openapi.SchemaRefs {
	p := openapi.SchemaRefs{}

	for i := 0; i < len(keyVals); i = i + 2 {
		p.Set(keyVals[i].(string), &openapi.SchemaRef{
			Value: keyVals[i+1].(*openapi.Schema),
		})
	}

	return p
}
