package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	wordsearchsystemgrpc "github.com/chrisjpalmer/word_search_system_grpc"
)

type KeyWordsRouteHandler struct {
	wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient
}

func NewKeyWordsRouteHandler(wordSearchSystemClient wordsearchsystemgrpc.WordSearchSystemClient) *KeyWordsRouteHandler {
	newKeyWordsRouteHandler := new(KeyWordsRouteHandler)
	newKeyWordsRouteHandler.wordSearchSystemClient = wordSearchSystemClient
	return newKeyWordsRouteHandler
}

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

	switch r.Method {
	case "GET":
		//Apply search term and search for word
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

type KeyWordsRouteHandlerGetInput struct {
}

type KeyWordsRouteHandlerGetOutput struct {
	KeyWords []string `json:"keywords"`
}

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
