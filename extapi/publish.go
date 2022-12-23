package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func handlePublish(w http.ResponseWriter, r *http.Request) {

	postToRecorder(r)

	response := getLast()

	encoder := json.NewEncoder(w)
	encoder.Encode(response)

}

func postToRecorder(r *http.Request) {
	url := RECORDER_ENDPOINT + "/pub"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, reqErr := http.NewRequest("POST", url, r.Body)
	if reqErr != nil {
		fmt.Errorf("%w", reqErr)
	}

	req.Header.Set("User-Agent", "extapi")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Errorf("%w", getErr)
	}

	if res.Body == nil {
		fmt.Errorf("no response")
	} else {
		defer res.Body.Close()
	}
}

func getLast() []map[string]interface{} {
	url := RECORDER_ENDPOINT + "/api/0/last"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		fmt.Errorf("%w", reqErr)
	}

	req.Header.Set("User-Agent", "extapi")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		fmt.Errorf("%w", getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Errorf("%w", readErr)
	}

	// get json as a map with interface{}
	resultList := []map[string]interface{}{}
	err := json.Unmarshal([]byte(body), &resultList)
	if err != nil {
		fmt.Errorf("extended api: fail to unmarshal json: %w", err)
	}
	//log.Printf("INFO: jsonMap, %s\n", jsonMap)

	// results list
	//resultList := jsonMap["results"].([]interface{})
	//log.Printf("INFO: resultList, %s\n", resultList)

	responses := []map[string]interface{}{}

	var type_location interface{} = "location"
	var type_card interface{} = "card"

	for i := 0; i < len(resultList); i++ {
		var resultMap map[string]interface{} = resultList[i]

		tid := resultMap["tid"]
		lat := resultMap["lat"]
		lon := resultMap["lon"]
		tst := resultMap["tst"]
		vel := resultMap["vel"]
		alt := resultMap["alt"]
		acc := resultMap["acc"]
		vac := resultMap["vac"]

		if tid != nil && lat != nil && lon != nil && tst != nil {
			newLocation := make(map[string]interface{})
			newLocation["_type"] = type_location
			newLocation["tid"] = tid
			newLocation["lat"] = lat
			newLocation["lon"] = lon
			newLocation["tst"] = tst
			newLocation["vel"] = vel
			newLocation["alt"] = alt
			newLocation["acc"] = acc
			newLocation["vac"] = vac
			responses = append(responses, newLocation)

			name := resultMap["name"]
			face := resultMap["face"]

			if name != nil && face != nil {
				newCard := make(map[string]interface{})
				newCard["_type"] = type_card
				newCard["tid"] = tid
				newCard["name"] = name
				newCard["face"] = face
				responses = append(responses, newCard)
			}

		}

	}

	return responses

}
