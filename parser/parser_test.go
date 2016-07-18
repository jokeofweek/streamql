package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/aybabtme/streamql/ast"
)

func TestParse(t *testing.T) {

	tests := []struct {
		input string
		want  *ast.FiltersStmt
	}{
		{
			input: ".",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{},
					}},
				},
			},
		},
		{
			input: ".,.",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{},
					}},
					{Selectors: []*ast.SelectorStmt{
						{},
					}},
				},
			},
		},
		{
			input: ".|.",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{},
						{},
					}},
				},
			},
		},
		{
			input: ".hello",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
						}},
					}},
				},
			},
		},
		{
			input: `.[].id`,
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Array: &ast.ArraySelectorStmt{
							Each: &ast.EachSelectorStmt{},
							Child: &ast.SelectorStmt{
								Object: &ast.ObjectSelectorStmt{
									Member: "id",
								},
							},
						}},
					}},
				},
			},
		},
		{
			input: ".hello | .bye",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
						}},
						{Object: &ast.ObjectSelectorStmt{
							Member: "bye",
						}},
					}},
				},
			},
		},
		{
			input: ".hello , .bye",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
						}},
					}},
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "bye",
						}},
					}},
				},
			},
		},
		{
			input: ".hello    ",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
						}},
					}},
				},
			},
		},
		{
			input: `.hello\ \ \ \ `,
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello    ",
						}},
					}},
				},
			},
		},
		{
			input: ".[]",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Array: &ast.ArraySelectorStmt{
							Each: &ast.EachSelectorStmt{},
						}},
					}},
				},
			},
		},
		{
			input: ".hello[]",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
							Child: &ast.SelectorStmt{
								Array: &ast.ArraySelectorStmt{
									Each: &ast.EachSelectorStmt{},
								},
							},
						}},
					}},
				},
			},
		},
		{
			input: ".hello[1]",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
							Child: &ast.SelectorStmt{
								Array: &ast.ArraySelectorStmt{
									Index: &ast.IndexSelectorStmt{
										Index: 1,
									},
								},
							},
						}},
					}},
				},
			},
		},
		{
			input: ".hello[1:2]",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
							Child: &ast.SelectorStmt{
								Array: &ast.ArraySelectorStmt{
									RangeEach: &ast.RangeEachSelectorStmt{
										From: 1, To: 2,
									},
								},
							},
						}},
					}},
				},
			},
		},
		{
			input: ".hello[1:2][42][].world",
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						{Object: &ast.ObjectSelectorStmt{
							Member: "hello",
							Child: &ast.SelectorStmt{
								Array: &ast.ArraySelectorStmt{
									RangeEach: &ast.RangeEachSelectorStmt{
										From: 1, To: 2,
									},
									Child: &ast.SelectorStmt{
										Array: &ast.ArraySelectorStmt{
											Index: &ast.IndexSelectorStmt{
												Index: 42,
											},
											Child: &ast.SelectorStmt{
												Array: &ast.ArraySelectorStmt{
													Each: &ast.EachSelectorStmt{},
													Child: &ast.SelectorStmt{
														Object: &ast.ObjectSelectorStmt{
															Member: "world",
														},
													},
												},
											},
										},
									},
								},
							},
						}},
					}},
				},
			},
		},

		{
			input: `.[].id1 | .[].id2, .[].id3`,
			want: &ast.FiltersStmt{
				Filters: []*ast.FilterStmt{
					{Selectors: []*ast.SelectorStmt{
						// first
						{Array: &ast.ArraySelectorStmt{
							Each: &ast.EachSelectorStmt{},
							Child: &ast.SelectorStmt{
								Object: &ast.ObjectSelectorStmt{
									Member: "id1",
								},
							},
						}},
						// | second
						{Array: &ast.ArraySelectorStmt{
							Each: &ast.EachSelectorStmt{},
							Child: &ast.SelectorStmt{
								Object: &ast.ObjectSelectorStmt{
									Member: "id2",
								},
							},
						}},
					}},
					// , other filter
					{Selectors: []*ast.SelectorStmt{
						{Array: &ast.ArraySelectorStmt{
							Each: &ast.EachSelectorStmt{},
							Child: &ast.SelectorStmt{
								Object: &ast.ObjectSelectorStmt{
									Member: "id3",
								},
							},
						}},
					}},
				},
			},
		},
	}

	for n, tt := range tests {
		t.Logf("test #%d, input %q", n, tt.input)

		got, err := NewParser(strings.NewReader(tt.input)).Parse()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if !reflect.DeepEqual(tt.want, got) {

			t.Errorf("want=%#v", tt.want)
			t.Errorf(" got=%#v", got)

			wantData, _ := json.MarshalIndent(tt.want, "", "  ")
			gotData, _ := json.MarshalIndent(got, "", "  ")

			t.Errorf("want=%#v", tt.want)
			t.Errorf(" got=%#v", got)

			t.Errorf("want=%s", string(wantData))
			t.Errorf(" got=%s", string(gotData))

		}
	}
}
