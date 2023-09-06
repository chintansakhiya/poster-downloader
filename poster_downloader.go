package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type getData struct {
	Results      []getPath `json:"results"`
	TotalResults int       `json:"total_results"`
}

type getPath struct {
	OriginalTitle string `json:"original_title"`
	Path          string `json:"poster_path"`
}
 
func Downloader(titles []string) {

	for i := 0; i < len(titles); i++ {

		pathString := `https://api.themoviedb.org/3/search/movie?api_key=15d2ea6d0dc1d476efbca3eba2b9bbfb&query=` + titles[i]
		res, err := http.Get(pathString)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()
		posterPath, _ := io.ReadAll(res.Body)

		var data getData
		json.Unmarshal(posterPath, &data)

		if data.TotalResults != 0 {

			fileNamae := strings.ReplaceAll(data.Results[0].Path, "/", "")
			err:=downloadFile(data.Results[0].Path, fileNamae)
			if err!=nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(titles[i] + " not found")
		}

	}

	//---------------------------------------------------------------------------
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get("https://image.tmdb.org/t/p/original" + URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create("img/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
