package ffmpeg

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
)

func Duration(filename string) string {
	cmd := exec.Command("ffmpeg", "-i", filename)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	if err != nil {
		log.Print(err)
	}

	durationRegex := regexp.MustCompile("Duration: (.*?),")

	result := durationRegex.FindStringSubmatch(buffer.String())
	if len(result) == 2 {
		return result[1]
	} else {
		return "ERROR getting duration"
	}

	//TODO not correct
	return ""
}