package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

type User struct {
	Key   *datastore.Key `datastore:"-"`
	ID    int64          `json:"id"`
	Email string         `json:"email"`
	Name  string         `json:"name"`
}

func GetUserByParentKey(c context.Context, pk *datastore.Key) (*User, error) {

	uu := []User{}
	keys, err := datastore.NewQuery("User").Ancestor(pk).Distinct().GetAll(c, &uu)
	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		uk := keys[0]
		user := uu[0]
		user.Key = uk
		user.ID = uk.IntID()
		return &user, nil
	}

	return nil, nil
}

func SaveUser(c context.Context, user *User) error {
	_, err := datastore.Put(c, user.Key, user)
	if err != nil {
		return err
	}

	return nil
}

func FindOrCreateUser(c context.Context) (*User, error) {
	u := user.Current(c)

	if u == nil {
		return nil, nil
	}

	pk := datastore.NewKey(c, "User", u.ID, 0, nil)
	user, err := GetUserByParentKey(c, pk)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	newUser := User{
		Email: u.Email,
	}

	k1 := datastore.NewIncompleteKey(c, "User", pk)
	if err != nil {
		return nil, err
	}

	k2, err := datastore.Put(c, k1, pk)
	if err != nil {
		return nil, err
	}

	newUser.Key = k2
	newUser.ID = k2.IntID()
	return &newUser, nil
}
