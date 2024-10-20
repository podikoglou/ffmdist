package main

import (
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("POST /", reqHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	// read form
	uploaded, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	name := header.Filename

	// create file
	tmpName := fmt.Sprintf("tmp-%s", name)

	// create temp file to store the video being sent
	file, err := os.Create(tmpName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// copy over data from form to file
	if _, err = io.Copy(file, uploaded); err != nil {
		panic(err)
	}

	// downgrade file quality
	input := ffmpeg.Input(tmpName)
	input = input.Output("pipe:1", ffmpeg.KwArgs{"format": "h264"})
	input = input.WithOutput(w)

	if err = input.Run(); err != nil {
		panic(err)
	}

	if err = os.Remove(tmpName); err != nil {
		panic(err)
	}

}
