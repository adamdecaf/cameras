package main

import (
	"io"
	"log"
	"net/http"
	mjpeg "github.com/marpie/go-mjpeg"
	opencv "github.com/lazywei/go-opencv/opencv"
)

type imageStream chan opencv.IplImage

func downloadFile(path string, out imageStream) {
	cap := opencv.NewFileCapture(path)
	if cap == nil {
		log.Printf("error reading video file at %s\n", path)
	}

	for {
		frame := cap.QueryFrame()
		if frame != nil {
			out <- *frame
		}
	}
	return
}

func downloadUrl(url string, out imageStream) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		close(out)
		return
	}

	for {
		img, err := mjpeg.Decode(response.Body)
		if err == io.EOF {
			close(out)
			return
		}
		if err != nil {
			log.Println(err)
			continue
		}
		if err == nil && img != nil {
			cvimg := opencv.FromImage(*img)
			if cvimg != nil {
				out <- *cvimg
			}
		}
	}
}
