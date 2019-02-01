package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
)

type WordsRouteHandler struct {
	wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient
}

func NewWordsRouteHandler(wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient) *WordsRouteHandler {
	newWordsRouteHandler := new(WordsRouteHandler)
	newWordsRouteHandler.wordSearchSystemClient = wordSearchSystemClient
	return newWordsRouteHandler
}

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

type WordsRouteHandlerGetInput struct {
	KeyWord string `json:"keyword"`
}

type WordsRouteHandlerGetOutput struct {
	Matches []string `json:"matches"`
}

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

type WordsRouteHandlerPostInput struct {
	Words []string `json:"words"`
}

type WordsRouteHandlerPostOutput struct {
}

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
