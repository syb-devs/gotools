package structinfo

import (
	"errors"
	"reflect"
)

// ErrInvalidType is raised when the type is for info extraction is not a struct
var ErrInvalidType = errors.New("type for decoding must be a struct or a pointer to one")

// Info stores information about a struct and its fields
type Info struct {
	NumFields int
	Fields    struct {
		Exported []FieldInfo
		All      []FieldInfo
	}
}

// FieldInfo holds info about a struct field
type FieldInfo struct {
	Index    int
	Name     string
	Kind     reflect.Kind
	Exported bool
	Tag      string
}

// Extract returns informatin about a struct extracted using reflection
func Extract(s interface{}, tagNamespace string) (*Info, error) {
	t := reflect.TypeOf(s)
	// dereference until not a pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}

	numFields := t.NumField()

	info := &Info{
		NumFields: numFields,
	}
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		fInfo := FieldInfo{
			Index:    i,
			Name:     f.Name,
			Kind:     f.Type.Kind(),
			Exported: isExported(f),
			Tag:      f.Tag.Get(tagNamespace),
		}
		info.Fields.All = append(info.Fields.All, fInfo)
		if fInfo.Exported {
			info.Fields.Exported = append(info.Fields.Exported, fInfo)
		}
	}
	return info, nil
}

// isExported  returns true if the struct field is exported.
func isExported(f reflect.StructField) bool {
	return len(f.PkgPath) == 0
}
