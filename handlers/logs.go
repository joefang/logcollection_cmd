package handlers

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
		log, err := getAllEvents(logFile)
		if err != nil {
			if err.Error() == "file is empty" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%q is empty", logFile)))
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("not able to show log for %q", logFile)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(log))
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
	numLastEvents, err := strconv.Atoi(lastNevents)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("last event is not valid"))
		return
	}

	log, err := getLastEvents(logFile, numLastEvents, false)
	if err != nil {
		if err.Error() == "exit status 1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("no events found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error getting events: %v", err)))
		}
		return
	}

	if filter != emptyString {
		log = getFilteredLog(log, filter)
		if log == emptyString {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("no events found"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(log))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(log))
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

func getFilteredLog(a string, filter string) string {
	logs := strings.Split(a, "\n")
	var newString string

	for _, v := range logs {
		if strings.Contains(v, filter) {
			newString = newString + v + "\n"
		}
	}
	return newString
}

func getLastEvents(logFile string, lastEvents int, getEverything bool) (string, error) {
	fileHandle, err := os.Open(filepath.Join("/var/log", logFile))
	if err != nil {
		return emptyString, errors.New("could not open file")
	}
	defer fileHandle.Close()

	line := emptyString
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	if filesize == 0 {
		return emptyString, errors.New("file is empty")
	}

	var lines string
	lineCount := 0

	for {
		cursor -= 1
		_, err := fileHandle.Seek(cursor, io.SeekEnd)
		if err != nil {
			return emptyString, err
		}

		char := make([]byte, 1)
		_, err = fileHandle.Read(char)
		if err != nil {
			return emptyString, err
		}

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // start a new line
			lineCount++
			lines = lines + line
			line = emptyString
		}
		line = fmt.Sprintf("%s%s", string(char), line) // build the line from chars

		if !getEverything {
			if lastEvents == lineCount {
				break
			}
		}

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}
	return lines, nil
}

func getAllEvents(logFile string) (string, error) {
	// number of last events will be ignored
	// when geteverything is true
	return getLastEvents(logFile, -1, true)
}
