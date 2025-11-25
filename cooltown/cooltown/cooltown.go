package cooltown

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

type outputTrack struct {
	Audio string
}

func findID(audioB64 string) (string, int64) {
	data := map[string]string{"Audio": audioB64}
	if val, ok := data["Audio"]; !ok || val == "" {
		return "", 0 // 400 - Audio not in data
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		return "", -2 // 500 - Could not marshal body
	}

	searchResp, err := http.Post("http://localhost:3001/search", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return "", -2 // 500 - Unable to access search API
	}

	var resMap map[string]interface{}
	err = json.NewDecoder(searchResp.Body).Decode(&resMap)
	if err != nil {
		return "", -2 // 500 - Could not decode body
	}
	songID, ok := resMap["Id"]
	if !ok || songID == "" {
		return "", 0 // 400 - songID not in data
	}

	url := "http://localhost:3000/tracks/" + songID.(string)
	tracksResp, err := http.Get(url)

	if tracksResp.StatusCode == 200 {
		if err != nil {
			return "", -2 // 500 - Unable to access tracks API
		}

		var res map[string]interface{}
		err = json.NewDecoder(tracksResp.Body).Decode(&res)
		if err != nil {
			return "", -2 // 500 - Could not decode body
		}
		songAud, ok := res["Audio"]
		if !ok || songAud == "" {
			return "", -1 // 404 - songAud not found in local db
		}

		audio := songAud.(string)
		return audio, 1
	} else if tracksResp.StatusCode == 404 {
		return "", -1 // 404 - songID not found in local db
	} else {
		return "", -2 // 500 - Internal Server Error
	}
}

func findFragment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400) // Empty Body - Bad Request
	}
	audioDict := map[string]interface{}{}
	if err := json.Unmarshal(body, &audioDict); err != nil {
		w.WriteHeader(500) // Unmarshal error - Internal Server Error
	}
	audioB64 := audioDict["Audio"].(string) // base64 encoded audio

	if audio, n := findID(audioB64); n > 0 {
		d := outputTrack{Audio: audio}
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
	r.HandleFunc("/cooltown", findFragment).Methods("POST")
	return r
}
