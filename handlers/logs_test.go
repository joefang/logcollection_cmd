package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var logSortedContent = "Jul  3 20:58:53 copper\nJul  3 20:58:40 weather\nJul  3 20:58:23 love\nJul  3 20:58:05 news\nJul  3 20:58:05 cars\nJul  3 20:57:40 sports\nJul  3 20:57:23 services\nJul  3 20:56:57 world\nJul  3 20:56:57 hello \n"

func Router(route string, f func(w http.ResponseWriter, r *http.Request)) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(route, f).Methods("GET")
	return router
}

func TestGetLogFile(t *testing.T) {
	type testCase struct {
		name               string
		expectedStatusCode int
		expectedBody       string
		testFile           string
	}
	testCases := []testCase{
		{
			name:               "happy path shows the entire content",
			expectedStatusCode: http.StatusOK,
			expectedBody:       logSortedContent,
			testFile:           "test.log",
		},
		{
			name:               "file found",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `"bogus.log" not found`,
			testFile:           "bogus.log",
		},
	}
	logDir = "../handlers/testfiles" // overriding log file directory
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", filepath.Join("/logs/files", test.testFile), nil)
			response := httptest.NewRecorder()
			Router("/logs/files/{file}", GetLogFile).ServeHTTP(response, request)
			assert.Equal(t, test.expectedStatusCode, response.Code)
			assert.Equal(t, test.expectedBody, response.Body.String(), "expected sorted content")
		})
	}
}

func TestGetLogEvents(t *testing.T) {
	type testCase struct {
		name               string
		lastevents         string
		expectedStatusCode int
		expectedBody       string
		testFile           string
	}
	testCases := []testCase{
		{
			name:               "expected 2 log events",
			lastevents:         "2",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Jul  3 20:58:53 copper\nJul  3 20:58:40 weather\n",
			testFile:           "test.log",
		},
		{
			name:               "expected 4 log events",
			lastevents:         "4",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Jul  3 20:58:53 copper\nJul  3 20:58:40 weather\nJul  3 20:58:23 love\nJul  3 20:58:05 cars\n",
			testFile:           "test.log",
		},
		{
			name:               "file found",
			lastevents:         "3",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `"bogus.log" not found`,
			testFile:           "bogus.log",
		},
	}
	logDir = "../handlers/testfiles" // overriding log file directory
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", filepath.Join("/logs/files", test.testFile, "lastevents", test.lastevents), nil)
			response := httptest.NewRecorder()

			Router("/logs/files/{file}/lastevents/{lastEvents:[0-9]+}", GetLogEvents).ServeHTTP(response, request)
			assert.Equal(t, test.expectedStatusCode, response.Code)
			assert.Equal(t, test.expectedBody, response.Body.String(), "expected sorted content")
		})
	}
}

func TestGetLogEvents_Filtered(t *testing.T) {
	type testCase struct {
		name               string
		lastevents         string
		expectedStatusCode int
		expectedBody       string
		testFile           string
		filter             string
	}
	testCases := []testCase{
		{
			name:               "expected 2 log events",
			lastevents:         "10",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "Jul  3 20:58:05 cars\n",
			testFile:           "test.log",
			filter:             "car",
		},
		{
			name:               "expected 4 log events",
			lastevents:         "10",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "no events found.",
			testFile:           "test.log",
			filter:             "apple",
		},
		{
			name:               "file found",
			lastevents:         "30",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `"bogus.log" not found`,
			testFile:           "bogus.log",
		},
	}
	logDir = "../handlers/testfiles" // overriding log file directory
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", filepath.Join("/logs/files", test.testFile, "lastevents", test.lastevents)+"?filter="+test.filter, nil)
			response := httptest.NewRecorder()

			Router("/logs/files/{file}/lastevents/{lastEvents:[0-9]+}", GetLogEvents).ServeHTTP(response, request)
			assert.Equal(t, test.expectedStatusCode, response.Code)
			assert.Equal(t, test.expectedBody, response.Body.String(), "expected sorted content")
		})
	}
}
