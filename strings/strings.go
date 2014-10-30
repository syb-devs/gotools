package strings

import (
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParseObjectId(s string) bson.ObjectId {
	if len(s) != 24 {
		var oid bson.ObjectId
		return oid
	}
	return bson.ObjectIdHex(s)
}

func ParseBool(s string) (bool, error) {
	if b, err := strconv.ParseBool(s); err == nil {
		return b, nil
	}

	if i, err := ParseInt(s); err == nil {
		if i == 0 {
			return false, nil
		} else {
			return true, nil
		}
	}

	switch strings.ToLower(s) {
	case "1", "true", "on", "yes", "enabled":
		return true, nil
	case "0", "false", "off", "no", "disabled", "":
		return false, nil
	default:
		return true, nil
	}
}

func ParseList(s string) ([]string, error) {
	return splitList(s, ",", false), nil
}

func ParseTags(s string) ([]string, error) {
	return splitList(s, ",", true), nil
}

func splitList(s string, sep string, lower bool) []string {
	frags := strings.Split(s, sep)
	for i, frag := range frags {
		frag = strings.Trim(frag, " ")
		if lower {
			frag = strings.ToLower(frag)
		}
		frags[i] = frag
	}
	return frags
}

func TakeWords(s string, num int) string {
	str := s
	words := strings.Split(str, " ")
	if len(words) > num {
		words = words[0:num]
	}
	return strings.Join(words, " ")
}

type List []string

func (sl List) Equals(list List) bool {
	if len(sl) != len(list) {
		return false
	}
	for i, item := range sl {
		if list[i] != item {
			return false
		}
	}
	return true
}

func (sl List) Contains(item string) bool {
	return sl.itemIndex(item) != -1
}

func (sl List) Remove(item string) List {
	i := sl.itemIndex(item)
	if i == -1 {
		return sl
	}
	var sl2 List
	sl2 = append(sl[:i], sl[i+1:]...)
	return sl2
}

func (sl List) itemIndex(item string) int {
	for i, s := range sl {
		if item == s {
			return i
		}
	}
	return -1
}
