package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func hello(w http.ResponseWriter, req *http.Request) {

	startTs := time.Now()

	//get the request method
	if req.Method == "POST" {

		// The request body just contain a string with a unique value
		// We can read the body of the request and print it
		readBody := make([]byte, req.ContentLength)
		err := error(nil)
		readBody, err = io.ReadAll(req.Body)
		if err != nil {
			reportTime(startTs, "Error reading request body", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s", readBody)))
		reportTime(startTs, fmt.Sprintf("%s", readBody), nil)

	} else {
		logrus.Errorf("Only POST please. Request method: %s", req.Method)
		w.Write([]byte("Only POST please"))
		http.Error(w, "Only POST please", http.StatusBadRequest)
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)
}

func reportTime(startTs time.Time, msg string, err error) {
	timeNow := time.Now()
	if err != nil {
		logrus.Errorf("Start:%v, Finish:%v, Elapsed(ms):%v Message: %s", startTs, timeNow, float64(timeNow.Sub(startTs).Nanoseconds())/1e6, err)
	} else {
		logrus.Infof("Start:%v, Finish:%v, Elapsed(ms):%v Message: %s", startTs, timeNow, float64(timeNow.Sub(startTs).Nanoseconds())/1e6, msg)
	}
}
