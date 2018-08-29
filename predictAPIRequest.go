package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"math"
)

//JSON structs for POST request
type ImageURL struct {
	Url string `json:"url"`
}

type Image struct {
	Image ImageURL `json:"image"`
}

type Data struct {
	Data Image `json:"data"`
}

var postMap map[string][]Data
var pMap = postMap

//buildJSON takes in an array of images and returns JSON for POST request
func buildJSON(line []string) []byte {
	var dataArray []Data
	for _, imgurl := range line {
		var urlMap = ImageURL{}
		urlMap.Url = imgurl
		var imageMap = Image{}
		imageMap.Image = urlMap
		var dataMap = Data{}
		dataMap.Data = imageMap
		dataArray = append(dataArray, dataMap)
	}
	pMap = make(map[string][]Data)
	pMap["inputs"] = dataArray
	pagesJson, err := json.Marshal(pMap)
	check(err)
	return pagesJson
}
var wg = sync.WaitGroup{}
//getPredictions divides the input images into chunks for predictions
func getPredictions(fileName string) {
	//lines,err := urlToLines("https://docs.google.com/document/d/13IOIc53stdGQEpR6ijamnhB04jG6P0ktln0JaazQTf4")
	lines, err := readLines(fileName)
	check(err)
	chunkSize := 128 //max size is 128
	numImages := 3//len(lines)
	var x = float64(numImages)/float64(chunkSize)
	var numRoutines = math.Ceil(x)
	wg.Add(int(numRoutines))
	for i := 0; i < numImages; i = i + chunkSize {
		if i+chunkSize > numImages {
			line := lines[i:numImages]
			go sendPredictRequest(line)
		}else {
			line := lines[i : i+chunkSize]
			go sendPredictRequest(line)
		}
	}
	wg.Wait()
}

//getModelID returns the modelID of the set model name
func getModelID(modelName string) string {
	if modelName == "General" {
		return "aaa03c23b3724a16a56b629203edc62c"
	}
	return "null"
}

//sendPredictRequest sends the post request for the chunk of images to the Clarifai PredictAPI
func sendPredictRequest(line []string) {
	var CLARIFAIPREDICTAPIKEY = getClarifAPIKey("Predict")
	var modelID = getModelID(modelName)
	var jsonStr = buildJSON(line)
	var reqURL = "https://api.clarifai.com/v2/models/" + modelID + "/outputs"
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Key "+CLARIFAIPREDICTAPIKEY)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	check(err)
	res, err1 := client.Do(req)
	check(err1)
	readResponse(res)
	wg.Done()
}
