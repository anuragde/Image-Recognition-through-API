package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var resultsMap map[string]map[string]float64
var topTagsMap map[string]map[int]map[string]float64
var fileName = "images.txt"
var fileUrl = "https://drive.google.com/uc?id=1wNolNaUqvDFAlfAYXOpsAjI0n1YGEMKFz95MIEd_mKc&export=download"
var modelName = "General"
var input string

//Check prints the error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//err := downloadFile("images.txt", fileUrl)
	//check(err)
	resultsMap = make(map[string]map[string]float64)
	topTagsMap = make(map[string]map[int]map[string]float64)
	fmt.Println("Please wait while predicting the input...")
	getPredictions(fileName)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter a tag: ")
		scanner.Scan()
		input = scanner.Text()
		getResults(strings.ToLower(input))
	}
}
