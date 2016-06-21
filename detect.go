package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path"
	"runtime"
	"time"
	opencv "github.com/lazywei/go-opencv/opencv"
)

// todo: new channel, in new file
// todo: only parse some frames, not all

func DetectFaces(in imageStream) {
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := opencv.LoadHaarClassifierCascade(path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml"))

	if cascade == nil {
		log.Println("error loading classifier")
		return
	}

	for {
		cvimg, ok := <-in
		if ok {
			ci := cvimg
			go detect(*cascade, ci)
		}
	}
}

func detect(cascade opencv.HaarCascade, img opencv.IplImage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovering from panic")
		}
	}()

	detectFaces(cascade, img)
}

func detectFaces(cascade opencv.HaarCascade, img opencv.IplImage) {
	faces := cascade.DetectObjects(&img)
	log.Printf("found %d faces", len(faces))
	for _, value := range faces {
		if (value.X() != nil && value.Y() != nil) && (value.Height() != nil && value.Width() != nil) {
			opencv.Rectangle(&img,
				opencv.Point{value.X() + value.Width(), value.Y()},
				opencv.Point{value.X(), value.Y() + value.Height()},
				opencv.ScalarAll(255.0), 1, 1, 0)
		}
	}

	if len(faces) > 0 {
		go write(img)
	}
}

func write(ci opencv.IplImage) {
	img := ci.ToImage()
	when := time.Now().UnixNano() / int64(time.Millisecond)

	f, _ := os.Create(fmt.Sprintf("face-%d.jpeg", when))
	defer f.Close()

	jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
}
