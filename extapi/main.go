package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
)

const RECORDER_ENDPOINT string = "http://localhost:8083"

var friendsStore *FileStore

func main() {

	fmt.Println("extended api: Starting")

	friendsStore, storeError := NewFileStore("/owntracks-storage/extapi-store/friends.json")
	if storeError != nil {
		fmt.Errorf("extended api: fatal error when setting up friendsStore: %w", storeError)
	}

	http.Handle("/0/friends", &handleFriends{friendsStore: friendsStore})
	http.HandleFunc("/0/publish", handlePublish)

	listener, listenError := net.Listen("tcp", ":491")
	if listenError != nil {
		fmt.Errorf("extended api: fatal error: %w", listenError)
	}
	defer listener.Close()

	serveError := fcgi.Serve(listener, nil)
	if serveError != nil {
		fmt.Errorf("extended api: fatal error: %w", serveError)
	}

}
