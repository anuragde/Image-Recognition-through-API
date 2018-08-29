package main

var CLARIFAIPREDICTAPIKEY = "" //Your API_Key

func getClarifAPIKey(keyType string) string {
	if keyType == "Predict" {
		return CLARIFAIPREDICTAPIKEY
	}
	return "null"
}
