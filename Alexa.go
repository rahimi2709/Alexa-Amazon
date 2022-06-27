package main

import (
	"bytes"
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	//"strings"
)

// To set all URI values as our running microservices(SST, Alpha And TTS)
const (
	URIAlpha = "http://localhost:3001/alpha"
	URISTT   = "http://localhost:3002/stt"
	URITTS   = "http://localhost:3003/tts"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
//Below function manage the transferring data between STT, Alpha and TTS microservices:
func alexa(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	check(err)
	reqBodybyte := []byte(reqBody)
	var PostedArgument map[string]string
	err = json.Unmarshal(reqBodybyte, &PostedArgument)
	check(err)
	TextOfSpeech, err2 := SpeechToText(string(PostedArgument["speech"]))
	check(err2)

	AnswerOfAlpha, err3 := Alpha(TextOfSpeech)
	check(err3)

	SpeechOfAnswer, err4 := TextToSpeech(AnswerOfAlpha)
	check(err4)

	w.Write([]byte(SpeechOfAnswer))
}

// Below 3 function are for transfering Requests and Responses between Alexa and other 3 microservices:
func SpeechToText(speech string) (string, error) {
	var jsonData = []byte(`{"speech":"` + speech + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URISTT, bytes.NewBuffer(jsonData))
	check(err)
	resp, er := client.Do(req)
	check(er)
	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	check(err3)
	return string(body), nil
}

func TextToSpeech(text string) (string, error) {
	var jsonData = []byte(`{"text":"` + text + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URITTS, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, er := client.Do(req)
	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	check(err3)
	check(er)
	check(err)
	return string(body), nil
}

func Alpha(text string) (string, error) {
	var jsonData = []byte(`{"text":"` + text + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URIAlpha, bytes.NewBuffer(jsonData))
	check(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return string(body), nil
}

// to keep running  the Alexa microservice in port 3000 to wait and received the posted data in POST method
//(Considering all other 3 Microservices are running simultaneously)

func main() {
	AlexaMicroService := mux.NewRouter()
	AlexaMicroService.HandleFunc("/alexa", alexa).Methods("POST")
	http.ListenAndServe(":3000", AlexaMicroService)
}
