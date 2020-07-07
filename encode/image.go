package encode

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Image encodes the image into a webp and returns the path to it
func Image(file multipart.File, user string, xOffset int, yOffset int, size int) error {

	workingDir, err := os.Getwd()

	// Create a temp file
	tempFile, err := ioutil.TempFile("files/temp-images", "image-*")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// read the uploaded file into a buffer and write it to our temp file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	tempFile.Write(fileBytes)

	// get the dimensions of the file that was uploaded and print to stdout
	infoCmd := exec.Command("identify", "-format", "%w %h", tempFile.Name())
	infoCmd.Dir = workingDir
	infoOutput, err := infoCmd.Output()
	if err != nil {
		return err
	}

	// If we got output without errors, then parse the output and log it.
	dimensionTokens := strings.Split(string(infoOutput), " ")
	width, err := strconv.Atoi(dimensionTokens[0])
	height, err := strconv.Atoi(dimensionTokens[1])
	if err != nil {
		return err
	}
	fmt.Printf("Image width: %v, height: %v\n", width, height)

	// Check if our crop fits within the image size
	if size+xOffset > width {
		return errors.New("Crop width exceeds the image dimensions")
	}

	if size+yOffset > height {
		return errors.New("Crop height exceeds the image dimensions")
	}

	// If the checks are good, let's crop our temp image.
	cropCmd := exec.Command("convert", tempFile.Name(), "-crop", fmt.Sprintf("%vx%v+%v+%v", size, size, xOffset, yOffset), "+repage", tempFile.Name())
	if err != nil {
		return err
	}
	cropCmd.Dir = workingDir
	var cropOutput bytes.Buffer
	cropCmd.Stderr = &cropOutput
	err = cropCmd.Run()
	if err != nil {
		return err
	}

	// since this is an image we'll use magick to encode it
	cmd := exec.Command("convert", tempFile.Name(), "(", "+clone", "-resize", "96x96^", "-write", fmt.Sprintf("files/thumbnails/%v.webp", user), "+delete", ")", "-resize", "512x512>", fmt.Sprintf("files/images/%v.webp", user))
	if err != nil {
		return err
	}
	cmd.Dir = workingDir
	var output bytes.Buffer
	cmd.Stderr = &output
	err = cmd.Run()
	fmt.Println(output.String())

	return err
}
