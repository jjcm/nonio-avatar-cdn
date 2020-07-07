package encode

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Image encodes the image into a webp and returns the path to it
func Image(file multipart.File, user string) error {

	workingDir, err := os.Getwd()

	// Create a temp file
	tempFile, err := ioutil.TempFile("files/temp-images", "image-*")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tempFile.Close()

	// read the uploaded file into a buffer and write it to our temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tempFile.Write(fileBytes)

	// get the dimensions of the file that was uploaded and print to stdout
	infoCmd := exec.Command("identify", "-format", "%w %h", tempFile.Name())
	infoCmd.Dir = workingDir
	infoOutput, err := infoCmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// If we got output without errors, then parse the output and log it.
	dimensionTokens := strings.Split(string(infoOutput), " ")
	width, err := strconv.Atoi(dimensionTokens[0])
	height, err := strconv.Atoi(dimensionTokens[1])
	if err != nil {
		fmt.Println("Error parsing image dimensions")
		return err
	}
	fmt.Printf("Image width: %v, height: %v\n", width, height)

	// since this is an image we'll use magick to encode it
	cmd := exec.Command("convert", tempFile.Name(), "(", "+clone", "-resize", "96x96^", "-write", fmt.Sprintf("files/thumbnails/%v.webp", user), "+delete", ")", "-resize", "512x512>", fmt.Sprintf("files/images/%v.webp", user))
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd.Dir = workingDir
	var output bytes.Buffer
	cmd.Stderr = &output
	err = cmd.Run()
	fmt.Println(output.String())

	return err
}
