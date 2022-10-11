package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"github.com/bogem/id3v2"
	"github.com/zackradisic/soundcloud-api"
	"strings"
)

func main() {
	csv := flag.String("csv", "songs.csv", "the csv file you want to read from")
	flag.Parse()
	_ = csv
	
	readCSV(*csv)

	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for index := range songs{
		url := songs[index].url

		if soundcloudapi.IsURL(url, true, true) {
			tracks, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
				URL: url,
			})
			if err != nil {
				log.Fatal(err.Error())
			}

			for index2 := range tracks {
				
				fileName, _ := cleanFileName(songs[index].title + "-" + songs[index].artist + ".mp3")
				
				if !fileExists(fileName){
					out, err := os.Create(fileName)

					if err != nil {
						fmt.Println(err.Error())
						return
					}
					defer out.Close()

					err = sc.DownloadTrack(tracks[index2].Media.Transcodings[0], out)
					
					if err != nil {
						log.Fatal(err.Error())
					}
					
					artworkURL := tracks[index2].ArtworkURL
					jpgName, _ := cleanFileName(songs[index].title + "-" + songs[index].artist + ".jpg")
					
					addArtwork(artworkURL, jpgName)
					
					fmt.Println(jpgName)

					tag, err := id3v2.Open(fileName, id3v2.Options{Parse: true})
					if err != nil {
						log.Fatal("Error while opening mp3 file: ", err)
					}
					defer tag.Close()

					artwork, err := os.ReadFile(jpgName)
					if err != nil {
						log.Fatal("Error while reading artwork file", err)
					}
				
					pic := id3v2.PictureFrame{
						Encoding:    id3v2.EncodingUTF8,
						MimeType:    "image/jpeg",
						PictureType: id3v2.PTFrontCover,
						Description: "Front cover",
						Picture:     artwork,
					}

					tag.AddAttachedPicture(pic)
					tag.SetTitle(songs[index].title)
					tag.SetArtist(songs[index].artist)

					//err = os.Remove("jpgName")
					//if err != nil {
					//	log.Fatal(err)
					//}

					fmt.Println("Downloaded " + fileName)
					
				}
			}
		} else {
			fmt.Println("URL #"+ strconv.Itoa(index + 1) +" is invalid")
		}
	}	
}

type song struct {
	title string
	artist string
	url string
}

var songs []song

func readCSV(fileName string) error{
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(io.Reader(csvFile))
	reader.FieldsPerRecord = 3

	for{
		line, error := reader.Read()
		
		if error == io.EOF {
			break
		} else if error != nil {
			return error
		}

		title := line[0]
		artist := line[1]
		url := line[2]
		
		songs = append(songs, song{title, artist, url})
	}
	return nil
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {			  
		return true
	} else {return false}
}

func addArtwork(url, fileName string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", err
	}	

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	
	return fileName, nil
}

func cleanFileName(fileName string) (string, error) {
	cleanName := strings.ReplaceAll(fileName, " ", "-")
	r := regexp.MustCompile(`[<>:"/\\|?*]`)
	cleanName = r.ReplaceAllString(cleanName, "")

	return cleanName, nil
}