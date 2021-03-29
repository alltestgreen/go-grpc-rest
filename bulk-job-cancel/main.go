package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	tuiDisputeCenter = "tuiDisputeCenter"
	clientToken      = ""
	clientID         = 1099
)

func main() {
	errFile, err := os.OpenFile("failed.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	e := executor{
		client:  &http.Client{},
		url:     "http://localhost:5219/api/monitoring/v2",
		failed:  make(map[string]string, 0),
		errFile: errFile,
	}

	clientKeys := []string{""}

	for _, clientKey := range clientKeys {
		e.cancelEnrollment(clientKey, tuiDisputeCenter, clientID)
	}

	fmt.Printf("Executed %d clientKeys, where %d was successful\n", len(clientKeys), e.success)
	if len(e.failed) > 0 {
		fmt.Printf("Failed cancellations: %+v", e.failed)
	}
}

type executor struct {
	client  *http.Client
	url     string
	success int
	failed  map[string]string
	errFile *os.File
}

func (e executor) cancelEnrollment(clientKey, enrollmentCode string, clientID int) {
	body := CancelMonitoringRequest{
		ClientKey:      clientKey,
		ClientID:       clientID,
		EnrollmentCode: enrollmentCode,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		e.handleError(err, clientKey)
	}

	// Create request
	req, err := http.NewRequest(http.MethodDelete, e.url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		e.handleError(err, clientKey)
		return
	}
	req.Header.Set("x-credmo-client-token", clientToken)
	req.Header.Set("x-credmo-client-id", strconv.Itoa(clientID))
	req.Header.Set("Content-Type", "application/json")

	// Do Request
	resp, err := e.client.Do(req)
	if err != nil {
		e.handleError(err, clientKey)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.handleError(err, clientKey)
		return
	}

	if resp.StatusCode == http.StatusOK {
		e.success++
	} else {
		e.handleError(fmt.Errorf("%s - %s", resp.Status, string(respBody)), clientKey)
	}
}

func (e executor) handleError(err error, clientKey string) {
	e.failed[clientKey] = err.Error()
	_, err = e.errFile.WriteString(clientKey + ", " + err.Error() + "\n")
}

type CancelMonitoringRequest struct {
	ClientKey      string `json:"clientKey"`
	ClientID       int    `json:"clientID"`
	EnrollmentCode string `json:"enrollmentCode"`
}
