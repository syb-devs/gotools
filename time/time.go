package time

import "time"

type Stamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ts *Stamps) Touch() {
	if ts.CreatedAt.IsZero() {
		ts.CreatedAt = time.Now()
	}
	ts.UpdatedAt = time.Now()
}

type DeleteStamps struct {
	Stamps
	DeletedAt time.Time `bson:",omitempty" json:"-"`
}

func (ts *DeleteStamps) MarkDeleted() {
	if ts.DeletedAt.IsZero() {
		ts.DeletedAt = time.Now()
	}
}
