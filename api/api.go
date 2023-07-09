package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"scraper/entity"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var Reg = regexp.MustCompile("[^a-zA-Z0-9-_]+")

func GetAlbums(userID string, apiKey string, album string) ([]entity.Album, error) {

	// Read the env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("The .env file is not loaded.")
	}

	// Get the environment variables
	API_BASE := os.Getenv("API_BASE")
	GET_ALBUMS := os.Getenv("GET_ALBUMS")
	FORMAT := os.Getenv("FORMAT")

	// URL
	url := API_BASE + "?method=" + GET_ALBUMS + "&api_key=" + apiKey + "&user_id=" + userID + "&format=" + FORMAT
	// fmt.Println(url)

	// GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var albums entity.Albums

	// Method 1:
	err = json.NewDecoder(resp.Body).Decode(&albums)
	if err != nil {
		return nil, err
	}

	// // Method 2:
	// body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	// // fmt.Println(string(body))              // convert to string before print
	// if err := json.Unmarshal([]byte(body), &albums); err != nil { // Parse []byte to go struct pointer
	// 	fmt.Println("Can not unmarshal JSON")
	// }
	// fmt.Println(albums.Stat)

	if albums.Stat != "ok" {
		return nil, errors.New("Flickr stated the status is not OK: " + albums.Stat)
	}

	// Only get one album if it's not all!
	if album != "all" {
		albumList := albums.Sets.Photoset // []Album
		for i := 0; i < int(albums.Sets.Total); i++ {
			currentAlbum := albumList[i]
			if currentAlbum.Title.Content == album {
				// Empty the array
				albumList = albumList[:0] // []
				// Only include one
				albumList = append(albumList, currentAlbum) //
				// fmt.Println(albumList)
				break
			}
		}
		return albumList, nil
	}

	return albums.Sets.Photoset, nil

}

func GetPhotosPerPage(userID string, apiKey string, photosetID string, page int) (entity.Photoset, error) {

	// Read the env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("The .env file is not loaded.")
	}

	// Get the environment variables
	API_BASE := os.Getenv("API_BASE")
	GET_PHOTOS := os.Getenv("GET_PHOTOS")
	FORMAT := os.Getenv("FORMAT")

	// URL
	url := API_BASE + "?method=" + GET_PHOTOS + "&api_key=" + apiKey + "&photoset_id=" + photosetID + "&user_id=" + userID + "&page=" + strconv.Itoa(page) + "&format=" + FORMAT
	// fmt.Println(url)

	// GET request
	resp, err := http.Get(url)
	if err != nil {
		return entity.Photoset{}, err
	}
	defer resp.Body.Close()

	var photos entity.Photos

	// Photos
	err = json.NewDecoder(resp.Body).Decode(&photos)
	if err != nil {
		return entity.Photoset{}, err
	}

	return photos.Photoset, nil

}

func GetPhotos(userID string, apiKey string, photosetID string) ([]entity.Photo, error) {

	// Variable to store all photos
	var allPhotos []entity.Photo

	// Call GetPhotosPerPage, for page 1
	currentPage := 1
	photos, err := GetPhotosPerPage(userID, apiKey, photosetID, currentPage)
	if err != nil {
		return []entity.Photo{}, err
	}

	allPhotos = append(allPhotos, photos.Photo...)

	// Check if it's more than 1 page
	if photos.Pages > 1 {
		for i := 2; i <= int(photos.Pages); i++ {
			currentPage = i
			photos, err := GetPhotosPerPage(userID, apiKey, photosetID, currentPage)
			if err != nil {
				return []entity.Photo{}, err
			}
			allPhotos = append(allPhotos, photos.Photo...)
		}
	}

	return allPhotos, nil

}

func GetPhotoSizeByPhoto(userID string, apiKey string, photoID string) ([]entity.Size, error) {

	// Read the env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("The .env file is not loaded.")
	}

	// Get the environment variables
	API_BASE := os.Getenv("API_BASE")
	GET_SIZES := os.Getenv("GET_SIZES")
	FORMAT := os.Getenv("FORMAT")

	// URL
	url := API_BASE + "?method=" + GET_SIZES + "&api_key=" + apiKey + "&photo_id=" + photoID + "&format=" + FORMAT
	// fmt.Println(url)

	// GET request
	resp, err := http.Get(url)
	if err != nil {
		return []entity.Size{}, err
	}
	defer resp.Body.Close()

	var links entity.Links

	// Links
	err = json.NewDecoder(resp.Body).Decode(&links)
	if err != nil {
		return []entity.Size{}, err
	}

	if links.Stat != "ok" {
		return nil, errors.New("Flickr stated the status is not OK: " + links.Stat)
	}

	return links.Sizes.Size, nil
}

func GetPhotoLinkByPhoto(userID string, apiKey string, photoID string) (string, error) {

	// Get the sizes
	sizes, err := GetPhotoSizeByPhoto(userID, apiKey, photoID)
	if err != nil {
		return "", err
	}

	// Get the photo link
	var result string
	for _, size := range sizes {
		if size.Label == "Original" {
			result = size.Source
			break
		}
	}

	return result, nil
}

func DownloadPhotoByLink(link string, filepath string, filename string) error {

	// Get the response bytes from the url
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Error response
	if resp.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	// Clean filename
	filename = Reg.ReplaceAllString(filename, "-")

	// Geth path + filename
	fileNameLength := 40
	if len(filename) > fileNameLength {
		filename = filename[:fileNameLength]
	}
	folder := path.Join(filepath, filename) + path.Ext(link)

	// Create a empty file
	file, err := os.Create(folder)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the downloaded bytes to the file created
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func DownloadAllPhotos(userID string, apiKey string, photos []entity.Photo, albumPath string) {

	// Channel to receive successful download notifications
	successChan := make(chan int)
	// Channel to receive failed download notifications
	failureChan := make(chan string)
	// Channel to notify that we're done!
	finishedChan := make(chan int)
	finishedDoneChan := make(chan int)
	// Channel for worker
	worker := make(chan int, 3)

	// WaitGroup to ensure all goroutines finish before exiting
	var wg sync.WaitGroup

	// Total data
	total := len(photos)
	wg.Add(total)

	// Start a Channel to track the progress of successful and failed downloads
	go func() {
		successCount := 0
		// failureCount := 0
		// var successCount, failureCount int
		for {
			select {
			case <-successChan:
				// wg.Done()
				successCount++
				fmt.Printf("\r Processing: %d out of %d photos", successCount, total)
			// case <-failureChan:
			// 	// wg.Done()
			// 	failureCount++
			case <-finishedChan:
				fmt.Printf("\n All downloads completed. Successful downloads: %d.\n", successCount)
				finishedDoneChan <- 1
				return
			}
		}
	}()

	// Spawn the Go functions from WaitGroup
	for i, photo := range photos {
		counter := strconv.Itoa(i + 1)
		counterTitle := counter + "-" + photo.Title
		// go DownloadPhotoAsync(&wg, worker, successChan, failureChan, userID, apiKey, albumPath, photo.ID, prefix+"-"+photo.Title)

		go func(photoID string, photoTitle string) {
			defer wg.Done()

			worker <- 1
			// Get the links
			link, err := GetPhotoLinkByPhoto(userID, apiKey, photoID)
			if err != nil {
				<-worker
				failureChan <- fmt.Sprintf("%s has error when getting link to download: %s", photoID, err)
				return
			}

			// Download and save
			err = DownloadPhotoByLink(link, albumPath, photoTitle)
			<-worker
			if err != nil {
				failureChan <- fmt.Sprintf("%s has error when downloading photo: %s", photoID, err)
				return
			}

			successChan <- 1
		}(photo.ID, counterTitle)

	}

	// Wait for all workers to finish
	wg.Wait()
	// time.Sleep(5 * time.Second) // a solution for data racing

	// Close the channels after all downloads are completed
	// close(successChan)
	close(failureChan)

	// Finished!
	finishedChan <- 1
	<-finishedDoneChan

}

// func DownloadPhotoAsync(wg *sync.WaitGroup, worker chan int, successChan chan<- int, failureChan chan<- string, userID string, apiKey string, albumPath string, photoID string, photoTitle string) {
// 	defer func() {
// 		wg.Done()
// 		// <-worker
// 	}()
// 	// defer <-worker

// 	worker <- 1
// 	// Get the links
// 	link, err := GetPhotoLinkByPhoto(userID, apiKey, photoID)
// 	if err != nil {
// 		<-worker
// 		failureChan <- fmt.Sprintf("%s has error when getting link to download: %s", photoID, err)
// 		return
// 	}

// 	// Download and save
// 	err = DownloadPhotoByLink(link, albumPath, photoTitle)
// 	<-worker
// 	if err != nil {
// 		failureChan <- fmt.Sprintf("%s has error when downloading photo: %s", photoID, err)
// 		return
// 	}

// 	successChan <- 1
// }
