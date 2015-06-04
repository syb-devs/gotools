package tsv_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/syb-devs/gotools/tsv"
)

type user struct {
	Name    string `tsv:"name"`
	Age     uint   `tsv:"age"`
	Range   int    `tsv:"range"`
	Active  bool   `tsv:",col:3"`
	Ignored string `tsv:"-"`
}

var decodeTests = []struct {
	in   string
	dest interface{}
	out  interface{}
	err  error
}{
	{
		in: `name	age	range	active
john doe	35	-10	1`,
		dest: &user{},
		out:  &user{"john doe", 35, -10, true, ""},
		err:  nil,
	},
	{
		in: `name	age	range	active
john doe	35	-10	false`,
		dest: &user{},
		out:  &user{"john doe", 35, -10, false, ""},
		err:  nil,
	},
	{
		in:   "",
		dest: &user{},
		out:  &user{},
		err:  io.EOF,
	},
}

func TestDecode(t *testing.T) {
	for i, test := range decodeTests {
		r := tsv.NewReader(strings.NewReader(test.in))
		r.ReadHeader()
		r.Next()
		err := r.Decode(test.dest)
		if err != test.err {
			t.Errorf("#%d: error mismatch\nhave %#+v\nwant %#+v", i, err, test.err)
		}
		if !reflect.DeepEqual(test.dest, test.out) {
			t.Errorf("#%d: mismatch\nhave %#+v\nwant %#+v", i, test.dest, test.out)
		}
	}
}
