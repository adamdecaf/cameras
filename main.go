package main

import (
	"flag"
	"log"
)

var (
	file string
	url string
)

func main() {
	flag.StringVar(&file, "file", "", "local filepath to read")
	flag.StringVar(&url, "url", "", "remote url to download")
	flag.Parse()

	stream := make(imageStream)
	if file != "" {
		go downloadFile(file, stream)
	}
	if url != "" {
		log.Printf("starting url download from %s", url)
		go downloadUrl(url, stream)
	}
	if file == "" && url == "" {
		log.Fatalf("no -file or -url params, one required")
	}

	// detect images off the stream
	DetectFaces(stream)
}
