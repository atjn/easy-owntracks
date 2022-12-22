package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-datastore"
	"io"
	"net/http"
)

func handleFriends(w http.ResponseWriter, r *http.Request) {

	user := getUser(r)

	if r.Method == "GET" {

		friendsList := getFriends(user)

		encoder := json.NewEncoder(w)
		encoder.Encode(friendsList)

	} else if r.Method == "POST" {

		friendsString, readError := io.ReadAll(r.Body)
		if readError != nil {
			fmt.Errorf("fatal")
		}

		friendsList := parseFriendsString(friendsString)

		setFriends(friendsList, user)

	}

}

func getFriends(user string) []string {

	hasUser := Has(user)

	if user != "" && hasUser {

		friendsString := Get(user)

		friendsList := parseFriendsString(friendsString)

		return friendsList

	} else {
		return []string{}
	}

}

func setFriends(friendsList []string, user string) {

	if user == "" {
		return
	}

	if len(friendsList) > 0 {

		friendsString, err := json.Marshal(friendsList)
		if err != nil {
			fmt.Errorf("extended api: failed to build database entry: %w", err)
		}

		Put(user, friendsString)

		fmt.Printf("extended api: User %s is now friends with: %s", user, friendsString)

	} else {

		Delete(user)

		fmt.Printf("extended api: User %s no longer has any friends.", user)

	}

}

func parseFriendsString(friendsString []byte) []string {

	friendsList := []string{}
	err := json.Unmarshal(friendsString, &friendsList)
	if err != nil {
		fmt.Errorf("extended api: failed to parse database entry: %w", err)
	}

	return friendsList

}

func getUser(r *http.Request) string {

	user := r.Header.Get("X-Limit-U")

	if user == "" {
		fmt.Errorf("extended api: fatal error: received a request to change friends, but no user was specified")
	}

	/*
		users := r.Header.Get("X-Limit-U")
		if len(users) == 1 {
			return users[0]
		} else if len(users) <= 0 {
			fmt.Errorf("extended api: fatal error: received a request to change friends, but no user was specified")
		} else {
			fmt.Errorf("extended api: fatal error: more than one user is not accepted")
		}
	*/

	return user

}

func Get(key string) []byte {
	userKey := datastore.NewKey(key)
	ctx := context.Background()
	friendsString, dberr := friendsDB.Get(ctx, userKey)
	if dberr != nil {
		fmt.Errorf("extended api: failed to read database entry: %w", dberr)
	}
	return friendsString
}

func Put(key string, value []byte) {
	userKey := datastore.NewKey(key)
	ctx := context.Background()
	dberr := friendsDB.Put(ctx, userKey, value)
	if dberr != nil {
		fmt.Errorf("extended api: failed to write database entry: %w", dberr)
	}
}

func Has(key string) bool {
	fmt.Printf("test: %s", key)
	userKey := datastore.NewKey(key)
	ctx := context.Background()
	doeshave, dberr := friendsDB.Has(ctx, userKey)
	if dberr != nil {
		fmt.Errorf("extended api: failed to read database entry: %w", dberr)
	}
	return doeshave
}

func Delete(key string) {
	userKey := datastore.NewKey(key)
	ctx := context.Background()
	dberr := friendsDB.Delete(ctx, userKey)
	if dberr != nil {
		fmt.Errorf("extended api: failed to remove database entry: %w", dberr)
	}
}
