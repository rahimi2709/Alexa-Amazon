package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func alpha(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	check(err)
	reqBodybyte := []byte(reqBody)
	var Question map[string]string
	err = json.Unmarshal(reqBodybyte, &Question)
	check(err)
	if Answer, err := alphaIN(Question["text"]); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(Answer))
	}
}

func alphaIN(text string) (string, error) {
	client := &http.Client{}

	// To send Question in appropriate format using GET method toward WolframAlpha and then receive the Answer from that:
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.wolframalpha.com/v1/result?i=%s&appid=T48R52-J4HRHP5U78", url.QueryEscape(text)), nil) // URI+text+"&appid="+appid, nil)
	check(err)
	rsp, err := client.Do(req)
	check(err)
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	check(err)
	return string(body), nil
}
// to keep running the microservice in port 3001 to wait and receive the posted data in POST method:

func main() {
	AlphaMicroService := mux.NewRouter()
	AlphaMicroService.HandleFunc("/alpha", alpha).Methods("POST")
	http.ListenAndServe(":3001", AlphaMicroService)
}
