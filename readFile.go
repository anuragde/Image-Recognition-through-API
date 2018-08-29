package main

import (
	"bufio"
	"log"
	"os"
)

//readLines reads the image url's from the file
func readLines(FileName string) ([]string, error) {
	file, err := os.Open(FileName)
	var lines []string
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err1 := scanner.Err(); err1 != nil {
		log.Fatal(err1)
	}
	return lines, nil
}

/*//readLinesReader reads lines from the passed reader
func readLinesReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

//urlToLines downloads the data from the URL
func urlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	return readLinesReader(resp.Body)
}
*/

/*//downloadFile downloads the file from the url and stores it in the local disk
func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	check(err)
	defer out.Close()

	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)

	return nil
}
*/
