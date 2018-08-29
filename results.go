package main

import "fmt"

//getResults prints the top 10 results for entered tag
func getResults(tag string) {
	var index = 1
	fmt.Println("The top images for " + tag + " are:")
	for index <= len(topTagsMap[tag]) {
		for key, _ := range topTagsMap[tag][index] {
			fmt.Println(key)
		}
		index = index + 1
	}
	fmt.Println("")
}