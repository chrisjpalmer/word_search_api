package main

import (
	"flag"
	"log"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
	"google.golang.org/grpc"
)

func main() {
	var (
		err    error
		config *Config
		conn   *grpc.ClientConn
	)
	//Get config
	configPath := flag.String("config", "config.json", "/path/to/config.json")
	flag.Parse()
	config, err = ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Announce start
	log.Println("WordSearchAPI has started")

	// Set up a connection to the word_search_system microservice.
	conn, err = grpc.Dial(config.WordSearchSystemAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	wordSearchSystemClient := wordsearchsystemgrpc.NewWordSearchSystemClient(conn)

	//Initialize the handlers and connect with HTTP
	wordsRouteHandler := NewWordsRouteHandler(wordSearchSystemClient)
	keyWordsRouteHandler := NewKeyWordsRouteHandler(wordSearchSystemClient)
	http.Handle("/words", wordsRouteHandler)
	http.Handle("/keywords", keyWordsRouteHandler)

	//Listen for http requests
	http.ListenAndServe(config.HTTPListenAddress, nil)
}
