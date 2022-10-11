package main

import(
	"fmt"
	"flag"
	"github.com/zackradisic/soundcloud-api"
	"log"
	"os"
	"encoding/csv"
	"io"
	"github.com/mikkyang/id3-go"
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
	fmt.Println(sc.ClientID())
	for index := range songs{
		url := songs[index].url
		fmt.Println(url)
		if soundcloudapi.IsURL(url, true, true) {
			tracks, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
				URL: url,
			})
			if err != nil {
				log.Fatal(err.Error())
			}

			for index2 := range tracks {
				fileName := songs[index].title + " - " + songs[index].artist + ".mp3"
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

					mp3File, err := id3.Open(fileName)

					if err != nil {
						log.Fatal(err.Error())
					}
					defer mp3File.Close()

					mp3File.SetArtist(songs[index].artist)
					mp3File.SetTitle(songs[index].title)
				}
			}
		} else {
			fmt.Println("URL is invalid")
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