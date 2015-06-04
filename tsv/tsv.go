package tsv

import (
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/syb-devs/gotools/structinfo"
)

var (
	// ErrEmptyRow is returned when there's no current row to decode
	ErrEmptyRow = errors.New("empty row")

	// ErrUnsuportedFieldType is returned when the underlying type of the struct field is not a supported one
	ErrUnsuportedFieldType = errors.New("struct field has an unsupported type to decode from TSV string")

	// ErrColNumMismatch is returned when the number of columns of the row does not match with the number of columns of the header
	ErrColNumMismatch = errors.New("struct field has an unsupported type to decode from TSV string")
)

// Reader reads TSV data
type Reader struct {
	r       *csv.Reader
	numCols int
	row     []string
	header  []string
	err     error
}

// NewReader returns a new TSV Reader that reads from r
func NewReader(r io.Reader) *Reader {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	return &Reader{
		r: cr,
	}
}

// Read reads one record from r.
func (r *Reader) Read() ([]string, error) {
	row, err := r.read()
	if err != nil {
		r.err = err
	}
	return row, err
}

func (r *Reader) read() ([]string, error) {
	row, err := r.r.Read()
	if err != nil {
		return row, err
	}

	if r.numCols > 0 && len(row) != r.numCols {
		return []string{}, ErrColNumMismatch
	}

	r.row = row
	r.numCols = len(row)
	return row, nil
}

// ReadHeader reads one row from the TSV file and stores its contents as the header definition
func (r *Reader) ReadHeader() ([]string, error) {
	h, err := r.Read()
	if err != nil {
		return h, err
	}
	r.header = h
	return h, nil
}

// Next reads the next row of the TSV file
func (r *Reader) Next() bool {
	_, err := r.Read()
	return err == nil
}

// Decode decodes the current row into dest struct
func (r *Reader) Decode(dest interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(r.row) == 0 {
		return ErrEmptyRow
	}

	info, err := structinfo.Extract(dest, "tsv")
	if err != nil {
		return err
	}
	v := reflect.ValueOf(dest)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for _, fi := range info.Fields.Exported {
		if fi.Tag == "-" {
			continue
		}
		f := v.Field(fi.Index)
		if f.CanSet() {
			col := r.getFieldCol(fi.Tag)
			if col >= 0 {
				if err = setVal(r.row[col], f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *Reader) getFieldCol(tag string) int {
	parts := strings.Split(tag, ",")
	if parts[0] != "" {
		return r.colIndex(parts[0])
	}
	if len(parts) < 2 {
		return -1
	}
	colInfo := parts[1]
	if colInfo == "" {
		return -1
	}
	ciParts := strings.Split(colInfo, ":")
	if len(ciParts) != 2 {
		return -1
	}
	i, err := strconv.Atoi(ciParts[1])
	if err != nil {
		return -1
	}
	return i
}

func (r *Reader) colIndex(col string) int {
	for i, name := range r.header {
		if col == name {
			return i
		}
	}
	return -1
}

func setVal(val string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		field.SetUint(uint64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)

	}
	return nil
}
