package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/user"

	"testing"
)

func TestUser(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatal(err)
	}

	defer inst.Close()

	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}

	testFindOrCreateWithoutAuth(t, appengine.NewContext(req))

	u := &user.User{Email: "email@test.com"}
	aetest.Login(u, req)

	testFindOrCreateWithAuth(t, appengine.NewContext(req))
}

func testFindOrCreateWithoutAuth(t *testing.T, ctx context.Context) {
	u, err := FindOrCreateUser(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if u != nil {
		t.Fatal("Expected nil user")
	}
}

func testFindOrCreateWithAuth(t *testing.T, ctx context.Context) {
	u2, err := FindOrCreateUser(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if u2 == nil {
		t.Fatal("User should not be nil")
	}

	user := user.Current(ctx)
	if u2.Email != user.Email {
		t.Fatal("Email doesn't match")
	}

	if u2.Key.Parent().StringID() != user.ID {
		t.Fatal("Parent ID doesn't match")
	}
}
