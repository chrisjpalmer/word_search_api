package main

import (
	"log"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	wordSearchSystemClient := wordsearchsystemgrpc.NewWordSearchSystemClient(conn)

	wordsRouteHandler := NewWordsRouteHandler(wordSearchSystemClient)
	keyWordsRouteHandler := NewKeyWordsRouteHandler(wordSearchSystemClient)

	http.Handle("/words", wordsRouteHandler)
	http.Handle("/keywords", keyWordsRouteHandler)
	http.ListenAndServe(":8080", nil)
}
