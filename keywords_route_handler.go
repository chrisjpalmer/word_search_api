package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
)

//KeyWordsRouteHandler - handles http requests to /keywords and interacts with WordSearchSystemClient
type KeyWordsRouteHandler struct {
	wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient
}

//NewKeyWordsRouteHandler - creates a new KeyWordsRouteHandler and initializes it with the WordSearchSystemClient
func NewKeyWordsRouteHandler(wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient) *KeyWordsRouteHandler {
	newKeyWordsRouteHandler := new(KeyWordsRouteHandler)
	newKeyWordsRouteHandler.wordSearchSystemClient = wordSearchSystemClient
	return newKeyWordsRouteHandler
}

//ServeHttp - handles an http request at /keywords
func (keyWordsRouteHandler *KeyWordsRouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		//get top 5 search keywords
		output, err = keyWordsRouteHandler.ServeGET(w, r)
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

//KeyWordsRouteHandlerGetInput - the input parameters for /keywords GET
type KeyWordsRouteHandlerGetInput struct {
}

//KeyWordsRouteHandlerGetOutput - the output parameters for /keywords GET
type KeyWordsRouteHandlerGetOutput struct {
	KeyWords []string `json:"keywords"`
}

//ServeGET - handles a GET request at /keywords
func (keyWordsRouteHander *KeyWordsRouteHandler) ServeGET(w http.ResponseWriter, r *http.Request) (output *KeyWordsRouteHandlerGetOutput, err error) {
	//Deserialize input
	var input KeyWordsRouteHandlerGetInput
	jsonDec := json.NewDecoder(r.Body)
	jsonDec.Decode(&input)

	//Form request
	top5SearchKeyWordsRequest := &wordsearchsystemgrpc.Top5SearchKeyWordsRequest{}

	//Send request
	reply, _err := keyWordsRouteHander.wordSearchSystemClient.Top5SearchKeyWords(r.Context(), top5SearchKeyWordsRequest)
	if _err != nil {
		return nil, _err
	}

	//Form output response
	_output := &KeyWordsRouteHandlerGetOutput{KeyWords: reply.Keywords}

	//Return output
	return _output, nil
}
