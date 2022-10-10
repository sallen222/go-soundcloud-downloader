package main

import(
	"fmt"
	"flag"
	"github.com/zackradisic/soundcloud-api"
	"log"
	"os"
)

func main() {
	url := flag.String("url", "", "the url for the song you want to download")
	flag.Parse()
	_ = url
	
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(sc.ClientID())

	if soundcloudapi.IsURL(*url, true, true) {
		tracks, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
			URL: *url,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		fileName := tracks[0].Title + ".mp3"

		out, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer out.Close()

		err = sc.DownloadTrack(tracks[0].Media.Transcodings[0], out)
		if err != nil {
			log.Fatal(err.Error())
		}	
	} else {
		fmt.Println("URL is invalid")
	}
}