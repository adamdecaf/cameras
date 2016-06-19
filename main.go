package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path"
	"runtime"
	"time"
	opencv "github.com/lazywei/go-opencv/opencv"
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
		log.Println("no -file or -url params, one required")
	}

	// load up decoder
	_, currentfile, _, _ := runtime.Caller(0)
	cascade := opencv.LoadHaarClassifierCascade(path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml"))

	// read each image off the channel
	// todo: new channel, in new file
	// todo: only parse some frames, not all
	for {
		cvimg, ok := <-stream
		if ok {
			go func() {
				ci := cvimg
				log.Println("got image from stream")
				faces := cascade.DetectObjects(&ci)
				log.Printf("found %d faces", len(faces))
				for _, value := range faces {
					fmt.Println(value.X())
					fmt.Println(value.Width())
					fmt.Println(value.Y())

					opencv.Rectangle(&ci,
						opencv.Point{value.X() + value.Width(), value.Y()},
						opencv.Point{value.X(), value.Y() + value.Height()},
						opencv.ScalarAll(255.0), 1, 1, 0)
				}

				if len(faces) > 0 {
					img := ci.ToImage()
					when := time.Now().UnixNano() / int64(time.Millisecond)
					f, _ := os.Create(fmt.Sprintf("face-%d.jpeg", when))
					jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
					f.Close()
					log.Println("saved image to disk")
				}
			}()
		}
	}
}
