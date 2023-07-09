package main

import (
	"errors"
	"log"
	"os"
	"path"
	"scraper/api"
	"scraper/entity"
)

func main() {

	// Credential
	// Input in your User ID
	userID := ""
	// Input in your API key:
	apiKey := ""

	// Get albums
	// Input in your album:
	album_name := ""
	albums, err := api.GetAlbums(userID, apiKey, album_name)
	if err != nil {
		log.Println(err)
	}

	// Get all photos in all the albums
	var allPhotos []entity.Photo
	for _, album := range albums {
		photos, err := api.GetPhotos(userID, apiKey, album.ID)
		if err != nil {
			log.Println(err)
		}
		allPhotos = append(allPhotos, photos...)
	}

	// Put your root path here
	rootPath := ""
	if _, err := os.Stat(rootPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(rootPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Clean the album name, create the album folder
	album_name = api.Reg.ReplaceAllString(album_name, "-")
	albumPath := path.Join(rootPath, album_name)
	if _, err := os.Stat(albumPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(albumPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Download all the photos!
	api.DownloadAllPhotos(userID, apiKey, allPhotos, albumPath)
}
