package strings_test

import (
	"testing"

	"bitbucket.org/syb-devs/gotools/strings"
	"gopkg.in/mgo.v2/bson"
)

var emptyObjectId bson.ObjectId

var parseObjectIdTests = []struct {
	input    string
	expected bson.ObjectId
}{
	{"dsdhjskdasdas", emptyObjectId},
	{"54490565f4444bc089dfeb09", bson.ObjectIdHex("54490565f4444bc089dfeb09")},
}

func TestParseObjectId(t *testing.T) {
	for _, test := range parseObjectIdTests {
		actual := strings.ParseObjectId(test.input)
		if actual != test.expected {
			t.Errorf("expecting ParseObjectId(%v) to be %v, got %v", test.input, test.expected, actual)
		}
	}
}

var parseIntTests = []struct {
	input     string
	expected  int
	withError bool
}{
	{"0", 0, false},
	{"00", 0, false},
	{"127.45", 0, true},
	{"009999", 9999, false},
	{"235,43", 0, true},
	{"43a", 0, true},
}

func TestParseInt(t *testing.T) {
	for _, test := range parseIntTests {
		actual, err := strings.ParseInt(test.input)
		if actual != test.expected {
			t.Errorf("expecting ParseInt(%v) to be %v, got %v", test.input, test.expected, actual)
		}

		if err == nil {
			if test.withError {
				t.Errorf("expecting an error parsing string %v as int, but got none", test.input)
			}
		} else {
			if !test.withError {
				t.Errorf("expecting no error parsing string %v as int, got %v", test.input, err)
			}
		}
	}
}

var parseBoolTests = []struct {
	input    string
	expected bool
}{
	{"", false},
	{"0", false},
	{"OFF", false},
	{"False", false},
	{"DisAbled", false},
	{"1", true},
	{"On", true},
	{"TRUE", true},
	{"enabled", true},
	{"\n", true},
	{"00", false},
	{"0032", true},
}

func TestParseBool(t *testing.T) {
	for _, test := range parseBoolTests {
		actual, _ := strings.ParseBool(test.input)
		if actual != test.expected {
			t.Errorf("expecting ParseBool(%v) to be %v, got %v", test.input, test.expected, actual)
		}
	}
}

var parseListTests = []struct {
	input    string
	expected []string
}{
	{"one, two, three", []string{"one", "two", "three"}},
	{"ONe,2,    \nThree   ", []string{"ONe", "2", "\nThree"}},
}

func TestParseList(t *testing.T) {
	for _, test := range parseListTests {
		actual, _ := strings.ParseList(test.input)
		if !listEquals(actual, test.expected) {
			t.Errorf("expecting ParseList(%v) to be %v, got %v", test.input, test.expected, actual)
		}
	}
}

var parseTagsTests = []struct {
	input    string
	expected []string
}{
	{"one, two, three", []string{"one", "two", "three"}},
	{"ONe,2,    \nThree   ", []string{"one", "2", "\nthree"}},
}

func TestParseTags(t *testing.T) {
	for _, test := range parseTagsTests {
		actual, _ := strings.ParseTags(test.input)
		if !listEquals(actual, test.expected) {
			t.Errorf("expecting ParseTags(%v) to be %v, got %v", test.input, test.expected, actual)
		}
	}
}

var takeWordsTests = []struct {
	input    string
	num      int
	expected string
}{
	{"We all lie to ourselves to be happy", 3, "We all lie"},
	{"You don't want the truth. You make up your own truth", 11, "You don't want the truth. You make up your own truth"},
	{"It's beer o'clock, and I'm buying", 8, "It's beer o'clock, and I'm buying"},
}

func TestTakeWords(t *testing.T) {
	for _, test := range takeWordsTests {
		actual := strings.TakeWords(test.input, test.num)
		if actual != test.expected {
			t.Errorf("expecting TakeWords(%v, %v) to be %v, got %v", test.input, test.num, test.expected, actual)
		}
	}
}

var listEqualsTests = []struct {
	listA    strings.List
	listB    strings.List
	expected bool
}{
	{strings.List{}, strings.List{}, true},
	{strings.List{"one", "two"}, strings.List{"one", "two"}, true},
	{strings.List{"1", "2", "3"}, strings.List{"1", "3"}, false},
	{strings.List{"1", "2", "3"}, strings.List{"1", "2", "4"}, false},
}

func TestListEquals(t *testing.T) {
	var equalStr string
	for _, test := range listEqualsTests {
		actual := test.listA.Equals(test.listB)
		if test.expected {
			equalStr = "equal to"
		} else {
			equalStr = "different from"
		}
		if actual != test.expected {
			t.Errorf("expecting list %v to be  %v %v", test.listA, equalStr, test.listB, actual)
		}
	}
}

var listContainsTests = []struct {
	list     strings.List
	item     string
	expected bool
}{
	{strings.List{}, "", false},
	{strings.List{"one", "two"}, "two", true},
	{strings.List{"1", "2", "3"}, "1", true},
}

func TestListContains(t *testing.T) {
	for _, test := range listContainsTests {
		actual := test.list.Contains(test.item)
		if actual != test.expected {
			t.Errorf("expecting list %v .Contains(%v) to be %v, got %v", test.list, test.item, test.expected, actual)
		}
	}
}

var listRemoveTests = []struct {
	list     strings.List
	item     string
	expected strings.List
}{
	{strings.List{"one", "two", "three"}, "one", strings.List{"two", "three"}},
	{strings.List{"one", "two", "three"}, "two", strings.List{"one", "three"}},
	{strings.List{"one", "two", "three"}, "three", strings.List{"one", "two"}},
	{strings.List{"one", "two", "three"}, "four", strings.List{"one", "two", "three"}},
}

func TestStringListRemove(t *testing.T) {
	for _, test := range listRemoveTests {
		actual := test.list.Remove(test.item)
		if !actual.Equals(test.expected) {
			t.Errorf("testing list.Remove(), removing %v from %v, expected %v, actual: %v", test.item, test.list, test.expected, actual)
		}
	}
}

func listEquals(a, b []string) bool {
	return strings.List(a).Equals(strings.List(b))
}
