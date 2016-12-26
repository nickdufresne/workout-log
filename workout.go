package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"time"
)

type Workout struct {
	Key         *datastore.Key `datastore:"-"`
	UserID      string         `json:"-"`
	Type        string         `json:"type"`
	Duration    int64          `json:"duration"`
	Distance    int            `json:"distance"`
	DistanceUOM string         `json:"distance_uom"`
	Details     string         `json:"details"`
	Date        time.Time      `json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func CreateWorkoutForUser(c context.Context, w *Workout, u *User) error {
	p := u.Key
	k := datastore.NewIncompleteKey(c, "Workout", p)

	k2, err := datastore.Put(c, k, w)
	if err != nil {
		return err
	}

	w.Key = k2
	return nil
}
