package structinfo_test

import (
	"reflect"
	"testing"

	"github.com/syb-devs/gotools/structinfo"
)

type user struct {
	Name       string `json:"first_name"`
	Age        int    `struct:"sage" json:"age"`
	Active     bool   `struct:"data" json:"act"`
	unexported int    `struct:"303"`
}

var extractTests = []struct {
	in  interface{}
	tag string
	err error
	out *structinfo.Info
}{
	{
		in:  &user{},
		tag: "json",
		err: nil,
		out: &structinfo.Info{
			NumFields: 4,
			Fields: struct {
				Exported []structinfo.FieldInfo
				All      []structinfo.FieldInfo
			}{
				Exported: []structinfo.FieldInfo{
					structinfo.FieldInfo{Index: 0, Name: "Name", Kind: reflect.String, Exported: true, Tag: "first_name"},
					structinfo.FieldInfo{Index: 1, Name: "Age", Kind: reflect.Int, Exported: true, Tag: "age"},
					structinfo.FieldInfo{Index: 2, Name: "Active", Kind: reflect.Bool, Exported: true, Tag: "act"},
				},
				All: []structinfo.FieldInfo{
					structinfo.FieldInfo{Index: 0, Name: "Name", Kind: reflect.String, Exported: true, Tag: "first_name"},
					structinfo.FieldInfo{Index: 1, Name: "Age", Kind: reflect.Int, Exported: true, Tag: "age"},
					structinfo.FieldInfo{Index: 2, Name: "Active", Kind: reflect.Bool, Exported: true, Tag: "act"},
					structinfo.FieldInfo{Index: 3, Name: "unexported", Kind: reflect.Int, Exported: false, Tag: ""},
				},
			},
		},
	},
	{
		in: &user{
			Name:       "John Doe",
			Age:        35,
			Active:     true,
			unexported: 909,
		},
		tag: "struct",
		err: nil,
		out: &structinfo.Info{
			NumFields: 4,
			Fields: struct {
				Exported []structinfo.FieldInfo
				All      []structinfo.FieldInfo
			}{
				Exported: []structinfo.FieldInfo{
					structinfo.FieldInfo{Index: 0, Name: "Name", Kind: reflect.String, Exported: true, Tag: ""},
					structinfo.FieldInfo{Index: 1, Name: "Age", Kind: reflect.Int, Exported: true, Tag: "sage"},
					structinfo.FieldInfo{Index: 2, Name: "Active", Kind: reflect.Bool, Exported: true, Tag: "data"},
				},
				All: []structinfo.FieldInfo{
					structinfo.FieldInfo{Index: 0, Name: "Name", Kind: reflect.String, Exported: true, Tag: ""},
					structinfo.FieldInfo{Index: 1, Name: "Age", Kind: reflect.Int, Exported: true, Tag: "sage"},
					structinfo.FieldInfo{Index: 2, Name: "Active", Kind: reflect.Bool, Exported: true, Tag: "data"},
					structinfo.FieldInfo{Index: 3, Name: "unexported", Kind: reflect.Int, Exported: false, Tag: "303"},
				},
			},
		},
	},
	{
		in:  new(int),
		tag: "json",
		err: structinfo.ErrInvalidType,
		out: nil,
	},
}

func TestExtract(t *testing.T) {
	for i, test := range extractTests {
		info, err := structinfo.Extract(test.in, test.tag)
		if !reflect.DeepEqual(info, test.out) {
			t.Errorf("#%d: mismatch\nexpecting:\t%#+v\ngot:\t\t%#+v", i, test.out, info)
		}
		if err != test.err {
			t.Errorf("#%d: error mismatch\nexpecting: %+v\ngot: %+v", i, test.err, err)
		}
	}
}
