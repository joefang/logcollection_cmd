package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	pipe "github.com/b4b4r07/go-pipe"
	"github.com/gorilla/mux"
)

const emptyString = ""

var logDir = "/var/log"

// GetLogFile checks whether the file specified exist under /var/log
func GetLogFile(w http.ResponseWriter, r *http.Request) {
	logFile := mux.Vars(r)["file"]
	if !findFile(logFile) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%q not found", logFile)))
		return
	} else {
		logFromCat, err := getLogEvents(emptyString, logFile, emptyString, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("not able to show log for %q", logFile)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(logFromCat))
		return
	}
}

// GetLogEvents http handler to get log events
func GetLogEvents(w http.ResponseWriter, r *http.Request) {
	logFile := mux.Vars(r)["file"]
	lastNevents := mux.Vars(r)["lastEvents"]
	filter := r.FormValue("filter")

	if !findFile(logFile) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%q not found", logFile)))
		return
	}

	logFromTail, err := getLogEvents(lastNevents, logFile, filter, false)
	if err != nil {
		if err.Error() == "exit status 1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("no events found."))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error getting events: %v", err)))
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(logFromTail))
}

func getLogEvents(lastNevents string, logFile string, filter string, showFile bool) (string, error) {
	catFileCommand := exec.Command("cat", filepath.Join(logDir, logFile))
	tailCommand := exec.Command("tail", "-n", lastNevents, filepath.Join(logDir, logFile))
	sortCommand := exec.Command("sort", "-r")
	filterCommand := exec.Command("grep", filter)
	var err error
	var b bytes.Buffer
	if showFile {
		err = pipe.Command(&b, catFileCommand, sortCommand)
	} else if filter != emptyString {
		err = pipe.Command(&b, tailCommand, sortCommand, filterCommand)
	} else {
		err = pipe.Command(&b, tailCommand, sortCommand)
	}
	if err != nil {
		return emptyString, err
	}
	return b.String(), nil
}

func findFile(searchFile string) bool {
	var listOfLogFiles []string
	fileInfos, err := ioutil.ReadDir(logDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			listOfLogFiles = append(listOfLogFiles, fileInfo.Name())
		}
	}
	for _, file := range listOfLogFiles {
		if file == searchFile {
			return true
		}
	}
	return false
}
