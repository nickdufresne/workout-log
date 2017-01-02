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

func FindRecentWorkoutsForUser(c context.Context, u *User) ([]*Workout, error) {
	workouts := []*Workout{}
	keys, err := datastore.NewQuery("Workout").Ancestor(u.Key).Order("-Date").Limit(20).GetAll(c, &workouts)
	if err != nil {
		return nil, err
	}

	for idx, _ := range workouts {
		workouts[idx].Key = keys[idx]
	}

	return workouts, nil
}

func CreateWorkoutForUser(c context.Context, w *Workout, u *User) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	p := u.Key
	k := datastore.NewIncompleteKey(c, "Workout", p)

	k2, err := datastore.Put(c, k, w)
	if err != nil {
		return err
	}

	w.Key = k2
	return nil
}

func FindWorkoutForUser(c context.Context, u *User, encK string) (*Workout, error) {
	k, err := datastore.DecodeKey(encK)
	if err != nil {
		return nil, err
	}

	wo := new(Workout)
	if err := datastore.Get(c, k, wo); err != nil {
		return nil, err
	}

	wo.Key = k
	return wo, nil
}
