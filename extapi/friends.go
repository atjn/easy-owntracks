package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type handleFriends struct {
	friendsStore *FileStore
}

func (handler *handleFriends) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := getUser(r)

	if r.Method == "GET" {

		friendsList := getFriends(handler.friendsStore, user)

		encoder := json.NewEncoder(w)
		encoder.Encode(friendsList)

	} else if r.Method == "POST" {

		friendsString, readError := io.ReadAll(r.Body)
		if readError != nil {
			fmt.Errorf("fatal")
		}

		friendsList := parseFriendsString(friendsString)

		setFriends(handler.friendsStore, friendsList, user)

	}

}

func getFriends(friendsStore *FileStore, user string) []string {

	friendsList, hasUser := friendsStore.Get(user)

	if user != "" && hasUser {

		return friendsList

	} else {
		return []string{}
	}

}

func setFriends(friendsStore *FileStore, friendsList []string, user string) {

	if user == "" {
		return
	}

	if len(friendsList) > 0 {

		err := friendsStore.Set(user, friendsList)
		if err != nil {
			fmt.Errorf("extended api: failed to save to database: %w", err)
		}

		fmt.Printf("extended api: User %s is now friends with: %s\n", user, friendsList)

	} else {

		err := friendsStore.Delete(user)
		if err != nil {
			fmt.Errorf("extended api: failed to delete from database: %w", err)
		}

		fmt.Printf("extended api: User %s no longer has any friends.\n", user)

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

	return user

}
