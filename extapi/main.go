package main

import (
	"fmt"
	//	"bufio"
	"github.com/ipfs/go-ds-leveldb"
	"net"
	"net/http"
	"net/http/fcgi"
)

var RECORDER_ENDPOINT string = "http://localhost:8083"

var friendsDB *leveldb.Datastore

func main() {

	fmt.Println("extended api: Starting")

	friendsDB, storeError := leveldb.NewDatastore("/owntracks-storage/extapi-store/friendsDB", nil)
	if storeError != nil {
		fmt.Errorf("extended api: fatal error: %w", storeError)
	}
	defer friendsDB.Close()

	http.HandleFunc("/0/friends", handleFriends)
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
