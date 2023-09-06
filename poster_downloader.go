// package main
package MoviePoster

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type getData struct {
	Results      []getPath `json:"results"`
	TotalResults int       `json:"total_results"`
}

type getPath struct {
	OriginalTitle string `json:"original_title"`
	Path          string `json:"poster_path"`
}

func Downloader(movieTitles []string) {
	createFolder()
	var wg sync.WaitGroup
	for i := 0; i < len(movieTitles); i++ {
		temp := strings.ReplaceAll(movieTitles[i], " ", "%20")
		reqUrl := `https://api.themoviedb.org/3/search/movie?api_key=15d2ea6d0dc1d476efbca3eba2b9bbfb&query=` + temp
		res, err := http.Get(reqUrl)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()
		posterPath, _ := io.ReadAll(res.Body)

		var data getData
		json.Unmarshal(posterPath, &data)
		if data.TotalResults != 0 {

			fileNamae := strings.ReplaceAll(data.Results[0].Path, "/", "")
			wg.Add(1)
			go downloadFile(data.Results[0].Path, fileNamae, &wg, movieTitles[i])

		} else {
			fmt.Println(movieTitles[i] + " not found")
		}
	}
	wg.Wait()
}

func downloadFile(URL, fileName string, wg *sync.WaitGroup, movieTitles string) {
	//Get the response bytes from the url
	response, err := http.Get("https://image.tmdb.org/t/p/original" + URL)
	if err != nil {
		wg.Done()
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		wg.Done()
		fmt.Println("Received non 200 response code")
		return
	}
	//Create a empty file
	file, err := os.Create("img/" + fileName)
	if err != nil {
		wg.Done()
		fmt.Println(err)
		return
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		wg.Done()
		fmt.Println(err)
		return
	}
	wg.Done()
	fmt.Println("Downloaded: " + movieTitles)
}

func createFolder() {
	folderPath := "./img"

	// Check if the folder already exists
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// The folder does not exist, so create it
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folder:", err)
		} else {
			fmt.Println("Folder created successfully.")
		}
	} else if err != nil {
		// Handle other errors, such as permission issues
		fmt.Println("Error checking folder:", err)
	} else {
		fmt.Printf("%v Folder already exists.\n", folderPath)
	}
}
