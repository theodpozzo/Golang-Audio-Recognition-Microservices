package search

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type outputTrack struct {
	Id string
}

func Search(audioB64 string) (string, int64) {
	APIkey, err := ioutil.ReadFile("./api_key.txt")
	data := url.Values{
		"audio":     {audioB64},
		"api_token": {string(APIkey)},
	}

	APIresp, err := http.PostForm("https://api.audd.io/", data)
	if err != nil {
		return "", -2 //500
	}
	respBody, err := io.ReadAll(APIresp.Body)
	if err != nil {
		return "", -2 //500
	}

	var respDict map[string]interface{}
	if err := json.Unmarshal(respBody, &respDict); err != nil {
		return "", -2 //500
	}

	if respDict["status"] != "success" {
		errDict := respDict["error"]
		errMap := errDict.(map[string]interface{})
		errCode := errMap["error_code"]
		fmt.Println(errDict)
		if errCode.(float64) == 300 {
			return "", -1 /* Incorrect audio base64 - Throw 404 Not Found */
		} else if errCode.(float64) == 900 {
			return "", 0 /* Incorrect API key - Throw 400 Bad Request */
		}
		return "", 0 /* Status != "success" - Throw 400 Bad Request */
	}

	resultDict := respDict["result"]
	resultMap := resultDict.(map[string]interface{})
	songName := resultMap["title"]
	songName = strings.Replace(songName.(string), " ", "+", -1)
	return songName.(string), 1
}

func searchTrack(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400) /* Bad Request */
	}

	audioDict := map[string]interface{}{}
	if err := json.Unmarshal(body, &audioDict); err != nil {
		w.WriteHeader(500) /* Internal Server Error */
	}

	audioB64 := audioDict["Audio"].(string) // base64 encoded audio

	if songName, n := Search(audioB64); n > 0 {
		d := outputTrack{Id: songName}
		w.WriteHeader(200) /* OK */
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(400) /* Bad Request */
	} else if n == -1 {
		w.WriteHeader(404) /* Not Found */
	} else {
		w.WriteHeader(500) /* Internal Server Error */
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	fmt.Println("Microservice running...\n")
	/* Search */
	r.HandleFunc("/search", searchTrack).Methods("POST")
	return r
}
