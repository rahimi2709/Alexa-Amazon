package main
import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"fmt"
)
const (
	REGION = "uksouth"
	URItts = "https://" + REGION + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEY = "19c1cb3c0aa848608fed5a5a8a23d640"
)
func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
func TextToSpeech(w http.ResponseWriter, r *http.Request) {

	// to read the posted text
	PostedTexts := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&PostedTexts); err == nil {
		if text, ok := PostedTexts["text"]; ok {
			tt, err2 := TextToSpeechIN(text)
			check(err2)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(base64.StdEncoding.EncodeToString(tt)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Function text to speech which use Microsoft Azure text-to-sppech service
//to convert a text as a string to an encoded Base64 WAV file:
func TextToSpeechIN(text string) ([]byte, error) {

	//To make value of voice tag in correct way:
	var XMLtool bytes.Buffer
	xml.EscapeText(&XMLtool, []byte(text))
	PostedXML := []byte("<speak version=\"1.0\" xml:lang=\"en-US\">\n<voice xml:lang=\"en-US\" name=\"en-US-JennyNeural\">" + XMLtool.String() + "</voice>\n</speak>")
	client := &http.Client{}
	req, err := http.NewRequest("POST", URItts, bytes.NewBuffer(PostedXML))
	check(err)

	// set the required headers to use for converting in function
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
	rsp, err2 := client.Do(req)
	check(err2)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)

		return body, nil
	} else {
		return nil, errors.New("cannot convert text to speech")
	}
}

// to keep running the microservice in port 3003 to wait and receive the posted data in POST method:
func main() {
	TTSMicroService := mux.NewRouter()
	TTSMicroService.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", TTSMicroService)
}
