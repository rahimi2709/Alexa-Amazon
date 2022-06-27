package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
	KEY = "19c1cb3c0aa848608fed5a5a8a23d640"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

// Function Speech to text which use Microsoft Azure speech-to-text service
//to convert an encoded Base64 WAV file to a text as a string:

func SpeechToText(w http.ResponseWriter, r *http.Request) {

	// to read the encoded posted file content:
	reqBody, err := ioutil.ReadAll(r.Body)
	check(err)
	reqBodyByte := []byte(reqBody)
	var PostedData map[string]string
	err = json.Unmarshal(reqBodyByte, &PostedData)
	check(err)
	rawDecodedText, _ := base64.StdEncoding.DecodeString(PostedData["speech"])
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, bytes.NewReader([]byte(rawDecodedText)))
	check(err)

	// set the required headers to use for converting in function:
	req.Header.Set("Content-Type", "audio/wav;codecs=audio/pcm;samplerate=16000")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
	rsp, err := client.Do(req)
	check(err)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		var microsoftOutput map[string]string
		if err := json.Unmarshal(body, &microsoftOutput); err != nil {
			w.Write([]byte(microsoftOutput["DisplayText"]))
		}
	} else {
		w.Write([]byte(""))
	}
}

// to keep running  the microservice in port 3002 to wait and receive the posted data in POST method:
func main() {
	STTMicroService := mux.NewRouter()
	STTMicroService.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", STTMicroService)
}
