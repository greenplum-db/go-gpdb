package download


import (
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"errors"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/ielizaga/piv-go-gpdb/core"
	"bytes"
)

// Implementing the new PivNet API token system
// The below function extract the token ( UAA )
func getToken(uaa_token string) (string, error) {

	log.Info("Getting the access token from the UAA token provided.")

	body := AuthBody{RefreshToken: uaa_token}
	b, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal API token request body: %s", err.Error())
	}

	// Placing request for access token.
	req, err := http.NewRequest("POST", RefreshToken, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to construct API token request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API token request failed: %s", err.Error())
	}

	defer resp.Body.Close()

	// Is it a success or not
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch API token - received status %v", resp.StatusCode)
	}

	// Let store it for the rest of the program
	err = json.NewDecoder(resp.Body).Decode(&AuthenicationResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode API token response: %s", err.Error())
	}

	return AuthenicationResponse.Token, nil

}

// Get the Json from the provided URL or download the file
// if requested.
func GetApi(method string, urlLink string, download bool, filename string, filesize int64) ([]byte, error) {

	// Get the request
	req, err := http.NewRequest(method, urlLink, nil)
	core.Fatal_handler(err)

	// Add Header to the Http Request
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Adding Token header.
	// Verify there is a token on config file
	if core.IsValueEmpty(core.EnvYAML.Download.ApiToken) {
		return []byte(""), errors.New("Cannot find value for \"API_TOKEN\", check \"config.yml\"")
	} else {
		 //Initial kickoff will not have any valid token
		 //so we will request for the token
		if AuthenicationResponse.Token == "" {
			req.Header.Set("Authorization", "Bearer " + core.EnvYAML.Download.ApiToken)
			new_refresh_token, err := getToken(core.EnvYAML.Download.ApiToken)
			if err != nil || new_refresh_token == "" {
				return []byte(""), fmt.Errorf("Authentication Failure: %s", err.Error())
			}
			return []byte(""), nil
		} else { // hey we have the access token, so lets move on.
			req.Header.Set("Authorization", "Bearer " + AuthenicationResponse.Token)
		}
	}

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {return []byte(""), err}

	// If the status code is not 200, then error out
	if resp.StatusCode != http.StatusOK {

		// If we get 451 even after accepting the EULA, then the user would manually
		// connect to the webpage and accept the EULA as its pivotal legal policy
		// to accept the EULA if you are downloading the product for the first time.
		if resp.StatusCode == 451 {

			// Print to why we got this error
			fmt.Println("\n\x1b[33;1mReceived Status code 451 when access the API, this means as per the pivotal " +
				"legal policy if you are attempting to download the product for the first time you are requested to " +
				"to manually accept the end user license agreement (only one time). please connect " +
				"to PivNet and accept the end user license agreement and then try again, as this step cannot be avoided. Click on the link " +
				"mentioned below to redirect you to website to accept the EULA \x1b[0m")


			// Read the error text and store it
			bodyText, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var f interface{}
			_ = json.Unmarshal(bodyText, &f)
			m := f.(map[string]interface{})

			// Obtain the URL that the user can access to accept the EULA
			for k, v := range m {
				if k == "message" {
					fmt.Println("\n\x1b[35;1m" + v.(string) + "\x1b[0m\n")
				}
			}

		}
		return []byte(""), errors.New("API ERROR: HTTP Status code expected ("+ strconv.Itoa(http.StatusOK) +") / received (" + strconv.Itoa(resp.StatusCode) + "), URL (" + urlLink + ")")
	}

	// Close the body once its done
	defer resp.Body.Close()

	// If its to download the software then download it
	if download {

		// The Size of the file
		size := filesize

		// Fully qualifies path
		filename = core.DownloadDir + filename

		// Create th file
		out, err := os.Create(filename)
		if err != nil {return []byte(""), err}

		// Initalize progress bar
		done := make(chan int64)
		go PrintDownloadPercent(done, filename, int64(size))
		defer out.Close()

		// Start Downloading
		n, err := io.Copy(out, resp.Body)
		if err != nil {return []byte(""), err}
		done <- n
	}

	// Read the json
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {return []byte(""), err}
	return bodyText, nil
}
