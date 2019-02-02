package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
)

//WordsRouteHandler - handles http requests to /words and interacts with WordSearchSystemClient
type WordsRouteHandler struct {
	wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient
}

//NewWordsRouteHandler - creates a new WordsRouteHandler and initializes it with the WordSearchSystemClient
func NewWordsRouteHandler(wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient) *WordsRouteHandler {
	newWordsRouteHandler := new(WordsRouteHandler)
	newWordsRouteHandler.wordSearchSystemClient = wordSearchSystemClient
	return newWordsRouteHandler
}

//ServeHttp - handles an http request at /words
func (wordsRouteHandler *WordsRouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		output      interface{}
		jsonEncoder *json.Encoder
	)

	//Handle the error should one occur in this defer function
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorString := fmt.Sprintf("500 - \"%s\"", err.Error())
			w.Write([]byte(errorString))
		}
	}()

	//Based on the method, direct to the appropriate sub handler
	switch r.Method {
	case "GET":
		//Apply search term and search for word
		output, err = wordsRouteHandler.ServeGET(w, r)
	case "POST":
		//Submit words to be added
		output, err = wordsRouteHandler.ServePOST(w, r)
	default:
		err = errors.New("Invalid HTTP method supplied")
	}
	if err != nil {
		return
	}

	//Encode the response as JSON
	jsonEncoder = json.NewEncoder(w)
	err = jsonEncoder.Encode(output)
	if err != nil {
		return //Here for consistency's sake
	}
}

//WordsRouteHandlerGetInput - the input parameters for /words GET
type WordsRouteHandlerGetInput struct {
	KeyWord string `json:"keyword"`
}

//WordsRouteHandlerGetOutput - the output parameters for /words GET
type WordsRouteHandlerGetOutput struct {
	Matches []string `json:"matches"`
}

//ServeGET - handles a GET request at /words - executes the SearchWord query on the WordSearchSystem to find possible matches for the provided keyword
func (wordsRouteHander *WordsRouteHandler) ServeGET(w http.ResponseWriter, r *http.Request) (output *WordsRouteHandlerGetOutput, err error) {
	//Deserialize input
	var input WordsRouteHandlerGetInput
	jsonDec := json.NewDecoder(r.Body)
	jsonDec.Decode(&input)

	//Form request
	searchWordRequest := &wordsearchsystemgrpc.SearchWordRequest{
		KeyWord: input.KeyWord,
	}

	//Send request
	reply, _err := wordsRouteHander.wordSearchSystemClient.SearchWord(r.Context(), searchWordRequest)
	if _err != nil {
		return nil, _err
	}

	//Form output response
	_output := &WordsRouteHandlerGetOutput{Matches: reply.Matches}

	//Return output
	return _output, nil
}

//WordsRouteHandlerPostInput - the input parameters for /words POST
type WordsRouteHandlerPostInput struct {
	Words []string `json:"words"`
}

//WordsRouteHandlerPostOutput - the output parameters for /words POST
type WordsRouteHandlerPostOutput struct {
}

//ServePOST - handles a POST request at /words - executes the AddWords query on the WordSearchSystem to add new words to the list
func (wordsRouteHander *WordsRouteHandler) ServePOST(w http.ResponseWriter, r *http.Request) (output *WordsRouteHandlerPostOutput, err error) {
	//Deserialize input
	var input WordsRouteHandlerPostInput
	jsonDec := json.NewDecoder(r.Body)
	jsonDec.Decode(&input)

	//Form request
	addWordsRequest := &wordsearchsystemgrpc.AddWordsRequest{
		Words: input.Words,
	}

	//Send request
	_, _err := wordsRouteHander.wordSearchSystemClient.AddWords(r.Context(), addWordsRequest)
	if _err != nil {
		return nil, _err
	}

	//Form output response
	_output := &WordsRouteHandlerPostOutput{}

	//Return output
	return _output, nil
}
