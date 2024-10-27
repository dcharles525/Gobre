package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"net/http"
	"math/rand"
)

func main(){

	http.HandleFunc(
		"/convert-file/{ogFileType}/{newFileType}", 
		handleConvertFile,
	)
	startServer()

}

func startServer(){
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleConvertFile(
	responseWriter http.ResponseWriter, 
	request *http.Request,
){
	ogFileType := request.PathValue("ogFileType")
	newFileType := request.PathValue("newFileType")

	if len(ogFileType) == 0 || len(newFileType) == 0{
		http.Error(responseWriter, "Path parameters not included", 400)
	}

	if request.Method == "POST" {
		randBytes := strconv.Itoa(rand.Int())

		filePtr, _ := os.OpenFile(
			"/tmp/" + randBytes, 
			os.O_WRONLY|os.O_CREATE, 
			0644,
		)

		bodyBytes, requestBodyError := io.ReadAll(request.Body)

		if requestBodyError != nil {
			http.Error(responseWriter, fmt.Sprint(requestBodyError), 500)
		}

		_, fileWriteError := filePtr.Write(bodyBytes)

		if fileWriteError != nil {
			http.Error(responseWriter, fmt.Sprint(fileWriteError), 500)
		}
		
		filePtr.Close()

		libreofficeCmdError := runLibreoffice(ogFileType, newFileType, randBytes)
		if libreofficeCmdError != nil {
			http.Error(responseWriter, fmt.Sprint(libreofficeCmdError), 500)
		}

		data, readFileError := os.ReadFile("/tmp/" + randBytes + "." + newFileType)

		if readFileError != nil {
			http.Error(responseWriter, fmt.Sprint(readFileError), 500)
		}

		responseWriter.Write(data)
	} else {
		http.Error(responseWriter, "Invalid request method.", 405)
	}
}

func runLibreoffice(
	ogFileType string,
	newFileType string, 
	randBytes string,
)error{
	var libreofficeOptions string
	if ogFileType == "pdf" {
		libreofficeOptions = "--infilter='writer_pdf_import'"
	}
	libreofficeCmd := exec.Command(
		"bash",
		"-c",
		"libreoffice " + libreofficeOptions +
		" -env:UserInstallation=file:///tmp/" + randBytes + "_lo " + 
		"--headless --convert-to " + newFileType + " " + randBytes,
	)

	libreofficeCmd.Dir="/tmp"
	return libreofficeCmd.Run()
}
