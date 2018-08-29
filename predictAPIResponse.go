package main

import (
	"time"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"sync"
)

//JSON structs to parse the POST request

type ClarifaiAPIResponse struct {
	Status  Status   `json:"status"`
	Outputs []Output `json:"outputs"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Input struct {
	ID   string `json:"id"`
	Data Image  `json:"data"`
}

type Concepts struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	AppId string  `json:"app_id"`
}

type Model struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	CreatedAt    time.Time    `json:"created_at"`
	AppId        string       `json:"app_id"`
	OutputInfo   OutputInfo   `json:"output_info"`
	ModelVersion ModelVersion `json:"model_version"`
	DisplayName  string       `json:"display_name"`
}
type OutputInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	TypeExt string `json:"type_ext"`
}

type ModelVersion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"`
}
type OutputData struct {
	Concept []Concepts `json:"concepts"`
}

type Output struct {
	ID        string     `json:"id"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	Model     Model      `json:"model"`
	Input     Input      `json:"input"`
	Data      OutputData `json:"data"`
}

//helper maps
var temp = make(map[string]float64)
var topTagTemp = make(map[string]float64)
var topTagTemp1 = make(map[string]float64)
var mutex = &sync.Mutex{}
var urlV = ""

//readResponse reads the data from http response
func readResponse(res *http.Response){
	body, err2 := ioutil.ReadAll(res.Body)
	clarifaiAPIResponse := ClarifaiAPIResponse{}
	check(err2)
	defer res.Body.Close()
	resErr := json.Unmarshal(body, &clarifaiAPIResponse)
	check(resErr)
	defer processJSON(clarifaiAPIResponse)
}

//processJSON processes the response JSON and stores the Predicted values
func processJSON(clarifaiAPIResponse ClarifaiAPIResponse) {
	for i, _ := range clarifaiAPIResponse.Outputs {
		var url = clarifaiAPIResponse.Outputs[i].Input.Data.Image.Url
		for j, _ := range clarifaiAPIResponse.Outputs[i].Data.Concept {
			model := clarifaiAPIResponse.Outputs[i].Data.Concept[j].Name
			value := clarifaiAPIResponse.Outputs[i].Data.Concept[j].Value
			mutex.Lock()
			insertResultMap(model, url, value)
			insertTopTagsMap(model, url, value)
			mutex.Unlock()
		}
	}
}

//Inserts the values into result map
func insertResultMap(model string, url string, value float64) {
	if len(resultsMap[model]) != 0 {
		temp = resultsMap[model]
	}
	temp[url] = value
	resultsMap[model] = temp
}

//Inserts the values into toptags map
func insertTopTagsMap(model string, url string, value float64) {
	var flag = false
	flag = checkTopTag(model, url, value, flag)
	if flag == false {
		var index = len(topTagsMap[model])
		if index == 0 {
			topTagsMap[model] = make(map[int]map[string]float64)
			topTagsMap[model][index+1] = make(map[string]float64)
			topTagsMap[model][index+1][url] = value
		} else {
			if index < 10 {
				topTagsMap[model][index+1] = make(map[string]float64)
				topTagsMap[model][index+1][url] = value
			}
		}
	}
}

//Checks if the image(url) belongs to the top 10 images (prediction value in top)
//and updates the remaining images order if necessary
func checkTopTag(model string, url string, value float64, flag bool) bool {
	var index = 1
	var EPSILON float64 = 0.000001
	for index <= len(topTagsMap[model]) {
		if index > 10 {
			break
		}
		for _, l := range topTagsMap[model][index] {
			if (value - l) > EPSILON && flag==false {
				 	topTagTemp = topTagsMap[model][index]
					topTagsMap[model][index] = make(map[string]float64)
					topTagsMap[model][index][url] = value
					flag = true
					continue
			}
			if flag == true {
				topTagTemp1 = topTagsMap[model][index]
				topTagsMap[model][index] = topTagTemp
				topTagTemp = topTagTemp1
			}
		}
		index++
	}
	return flag
}
