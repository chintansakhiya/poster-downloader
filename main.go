package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type getData struct {
	Results      []getPath `json:"results"`
	TotalResults int       `json:"total_results"`
}

type getMovieName struct {
	Results []getPath `json:"data"`
}

type getPath struct {
	OriginalTitle string `json:"original_title"`
	Path          string `json:"poster_path"`
}

func main() {

	movieTitle, err := http.Get(`http://192.168.1.7:3000/api/v1/movies?limit=44932`)
	if err != nil {
		fmt.Println(err)
	}

	defer movieTitle.Body.Close()
	title, _ := ioutil.ReadAll(movieTitle.Body)

	var titleStr getMovieName
	json.Unmarshal(title, &titleStr)
	//----------------------------------------------------------------
	for i := 0; i < len(titleStr.Results); i++ {
		titleForPath := strings.Trim(titleStr.Results[i].OriginalTitle, " ")
		titleForPath = strings.ReplaceAll(titleForPath, " ", "%20")
		pathString := `https://api.themoviedb.org/3/search/movie?api_key=15d2ea6d0dc1d476efbca3eba2b9bbfb&query=` + titleForPath

		res, err := http.Get(pathString)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()
		posterPath, _ := ioutil.ReadAll(res.Body)

		var data getData
		json.Unmarshal(posterPath, &data)

		if data.TotalResults != 0 {

			fileNamae := strings.ReplaceAll(data.Results[0].Path, "/", "")
			downloadFile(data.Results[0].Path, fileNamae)
		} else {
			fmt.Println(titleForPath)

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
