package form

import (
	"net/http"
	"net/url"

	"github.com/syb-devs/gotools/strings"
	"gopkg.in/mgo.v2/bson"
)

type Form struct {
	values url.Values
}

func NewForm(r http.Request) (*Form, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	return &Form{values: r.Form}, nil
}

func (f Form) GetOne(field string) string {
	vals := f.values[field]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (f Form) Values() url.Values {
	return f.values
}

func (f Form) ParseTags(field string) ([]string, error) {
	tags := f.GetOne(field)
	return strings.ParseTags(tags)
}

func (f Form) ParseObjectId(field string) bson.ObjectId {
	id := f.GetOne(field)
	return strings.ParseObjectId(id)
}

func (f Form) ParseString(field string) string {
	return f.GetOne(field)
}

func (f Form) ParseBool(field string) (bool, error) {
	val := f.GetOne(field)
	return strings.ParseBool(val)
}
