package room

import (
	"github.com/golang/protobuf/proto"
)

func mustProtoMarshal(v proto.Message) []byte {
	data, err := proto.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

var sampleCompiledTree1 = Tree{
	Sections: []*Section{
		{
			Common: &Common{
				Type:    "standard",
				IdAttr:  "stuff",
				Classes: []string{"a", "b"},
			},
			Rows: []*Row{
				{
					Common: &Common{
						Type: "halfhalf",
					},
					Columns: []*Column{
						{
							Modules: []*Module{
								{
									Type: "text",
									Data: []byte(""),
									Dyn:  1,
								},
							},
						},
					},
				},
				{
					Common: &Common{
						Type: "thirdthirdthird",
					},
					Columns: []*Column{
						{
							Modules: []*Module{
								{
									Type: "html",
									Data: []byte("<h2>Here is stuff</h2>"),
								},
							},
						},
					},
				},
			},
		},
		{
			Common: &Common{
				Type:    "fullwidth",
				IdAttr:  "stuff",
				Classes: []string{"a", "b"},
			},
			Rows: []*Row{
				{
					Common: &Common{
						Type: "halfhalf",
					},
					Columns: []*Column{
						{
							Modules: []*Module{
								{
									Type: "nav",
									//Data: mustProtoMarshal(&Nav{
									//	Links: []*NavLink{},
									//}),
								},
							},
						},
					},
				},
				{
					Common: &Common{
						Type: "thirdtwothirds",
					},
				},
			},
		},
	},
}
