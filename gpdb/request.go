package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AuthBody struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResp struct {
	Token string `json:"access_token"`
}

// Implementing the new PivNet API token system
// The below function extract the token ( UAA )
func getToken() string {

	Info("Getting the access token from the UAA token provided.")

	body := AuthBody{RefreshToken: Config.DOWNLOAD.APITOKEN}
	b, err := json.Marshal(body)
	if err != nil {
		Fatalf("Failed to marshal API token request body: %s", err.Error())
	}

	// Placing request for access token.
	req, err := http.NewRequest("POST", RefreshToken, bytes.NewReader(b))
	if err != nil {
		Fatalf("Failed to construct API token request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Fatalf("API token request failed: %s", err.Error())
	}

	defer resp.Body.Close()

	// Is it a success or not
	if resp.StatusCode != http.StatusOK {
		Fatalf("Failed to fetch API token - received status %v", resp.StatusCode)
	}

	// Let store it for the rest of the program
	authenicationResponse := new(AuthResp)
	err = json.NewDecoder(resp.Body).Decode(&authenicationResponse)
	if err != nil {
		Fatalf("Failed to decode API token response: %s", err.Error())
	}

	return authenicationResponse.Token
}

// Generate the URL headers
func generateHandler(method, url, token string, download bool) (*http.Response) {
	// Create new http request
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		Fatalf("Encountered error when sending new request: %v", err)
	}

	// copy headers
	// Add Header to the Http Request
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	// Skip SSL stuffs
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	// perform request
	var client = new(http.Client)
	if download { // download can take longer, so no need for timeout here
		client = &http.Client{Transport: transport}
	} else { // set 30 second timeout for any http request
		client = &http.Client{Transport: transport, Timeout: 30 * time.Second}
	}
	resp, err := client.Do(request)
	if err != nil {
		Fatalf("Encountered error when requesting the data from http: %v", err)
	}

	return resp
}

// The below functions fetch data from the URL
func fetch(method, url, token string) ([]byte) {

	Debugf("Requesting data from the url: %s", url)
	var contents []byte

	// Get the responses from the url
	resp := generateHandler(method, url, token, false)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Fatalf("Encountered error when reading the data from http: %v", err)
	}

	// Check if the response is 200
	if resp.StatusCode != http.StatusOK {
		Fatalf("Encountered invalid status code from http: %v", resp.StatusCode)
	}

	return contents
}

// the below function does the get from the URL
func get(url, token string) ([]byte) {
	return fetch("GET", url, token)
}

// the below function does the get from the URL
func post(url, token string) ([]byte) {
	return fetch("POST", url, token)
}

// Download the file requested
func downloadProduct(url, token string, r Responses) {

	// The Size of the file
	size := r.UserRequest.ProductFileSize

	// Fully qualifies path
	r.UserRequest.ProductFileName = Config.DOWNLOAD.DOWNLOADDIR + r.UserRequest.ProductFileName

	// Check if the file already exists. Skip download if the file is present
	// Only would work with gpdb, for other flags like gpcc / gpextra this won't work
	// since we don't store or request for any version information
	if cmdOptions.Product == "gpdb" {
		filePath, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("*%s*.zip", cmdOptions.Version))
		if len(filePath) > 0 && !cmdOptions.Always {
			Warnf("File %s found. Skipping download", filePath[0])
			Warn("To force re-download of the file, use -a flag")
			return
		}
	}

	// Create th file
	out, err := os.Create(r.UserRequest.ProductFileName)
	if err != nil {
		Fatalf("Creation of the file %s failed, err: %v", r.UserRequest.ProductFileName, err)
	}

	response := generateHandler("GET", url, token, true)
	defer response.Body.Close()

	// Initalize progress bar
	done := make(chan int64)
	go PrintDownloadPercent(done, r.UserRequest.ProductFileName, int64(size))
	defer out.Close()

	// Start Downloading
	n, err := io.Copy(out, response.Body)
	if err != nil {
		Fatalf("Error in downloading the file: %v", err)
	}
	done <- n
}
