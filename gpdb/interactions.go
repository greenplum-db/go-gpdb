package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Ask user what is the choice from the list provided.
func PromptChoice(TotalOptions int) int {
	Debugf("Promoting for choice from the total: %d", TotalOptions)

	var choiceEntered int
	fmt.Print("\nEnter your choice from the above list (eg.s 1 or 2 etc): ")

	// For unit test cases to work
	if !IsValueEmpty(os.Getenv("TEST_PROMPT_CHOICE")) {
		choiceEntered = strToInt(os.Getenv("TEST_PROMPT_CHOICE"))
	} else {
		// Start the new scanner to get the user input
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {

			// The choice entered
			choiceEntered, err := strconv.Atoi(input.Text())

			// If user enters a string instead of a integer then ask to re-enter
			if err != nil {
				fmt.Println("Incorrect value: Please choose a integer (eg.s 1 or 2 etc) from the above list")
				return PromptChoice(TotalOptions)
			}

			// If its a valid value move on
			if choiceEntered > 0 && choiceEntered <= TotalOptions {
				return choiceEntered
			} else { // Else ask for re-entering the selection
				fmt.Println("Invalid Choice: The choice you entered is not on the list above, try again.")
				return PromptChoice(TotalOptions)
			}
		}
	}
	return choiceEntered
}

// Prompt for confirmation
func YesOrNoConfirmation() string {
	Debugf("Promoting for yes or no confirmation")

	var YesOrNo = map[string]string{"y": "y", "ye": "y", "yes": "y", "n": "n", "no": "n"}

	// For unit test cases to work
	if !IsValueEmpty(os.Getenv("TEST_YES_CONFIRMATION")) {
		return os.Getenv("TEST_YES_CONFIRMATION")
	} else {
		// Start the new scanner to get the user input
		fmt.Print("You can use \"gpdb env -v <version>\" to set the env, do you wish to continue (Yy/Nn)?: ")
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {

			// The choice entered
			choiceEntered := input.Text()

			// If its a valid value move on
			if YesOrNo[strings.ToLower(choiceEntered)] == "y" { // Is it Yes
				return choiceEntered
			} else if YesOrNo[strings.ToLower(choiceEntered)] == "n" { // Is it No
				return choiceEntered
			} else { // Invalid choice, ask to re-enter
				fmt.Println("Invalid Choice: Please enter Yy/Nn, try again.")
				return YesOrNoConfirmation()
			}
		}
	}
	return ""
}

// provide choice of which version to download
func (r *Responses) ShowAvailableVersion(token string) {
	Info("Checking for the available version")

	// Local storehouse
	var ReleaseOutputMap = make(map[string]string)
	var ReleaseVersion []string
	var ReleaseFileOutPutMap = make(map[string]string)
	var ReleaseSizeOutPutMap = make(map[string]int64)

	// Get all the releases from the ReleaseJson
	for _, release := range r.ReleaseList.Release {
		ReleaseOutputMap[release.Version] = release.Links.Self.Href
		ReleaseVersion = append(ReleaseVersion, release.Version)
	}

	// Get the user choice
	_, r.UserRequest.releaseLink, _ = pickUserChoice(
		ReleaseOutputMap, ReleaseFileOutPutMap, ReleaseVersion, ReleaseSizeOutPutMap)
	r.UserRequest.versionChoosen = cmdOptions.Version

	// Since we have now extracted the version, now get the download URL
	r.ExtractDownloadURL(token)
}

// Extract all the tags of the open source github release
func ShowOpenSourceAvailableVersion(g GithubReleases) (string, string, int64) {
	Debugf("Showing all the available open source version available for download")

	// Local storehouse
	var ReleaseOutputMap = make(map[string]string)
	var ReleaseVersion []string
	var ReleaseFileOutPutMap = make(map[string]string)
	var ReleaseSizeOutPutMap = make(map[string]int64)

	// Regex compile for open source
	rx, _ := regexp.Compile("(?i)" + rx_open_source_gpdb)
	if getSystemInfoAndCheckIfItsUbuntu() {
		Debug("Changed the regex for finding the ubuntu product in open source")
		rx, _ = regexp.Compile("(?i)" + rx_gpdb_ubuntu_release_open_source)
	}

	// Get all the releases from the ReleaseJson
	for _, release := range g {
		for _, asset := range release.Assets {
			if rx.MatchString(asset.Name)  {
				tag := release.TagName
				ReleaseOutputMap[tag] = asset.BrowserDownloadURL
				ReleaseFileOutPutMap[tag] = asset.Name
				ReleaseSizeOutPutMap[tag] = asset.Size
				ReleaseVersion = append(ReleaseVersion, release.TagName)
			}
		}
	}
	return  pickUserChoice(
		ReleaseOutputMap, ReleaseFileOutPutMap, ReleaseVersion, ReleaseSizeOutPutMap)
}

// Pick the version from the user choice
func pickUserChoice(ReleaseOutputMap, ReleaseFileOutPutMap map[string]string,
	ReleaseVersion []string, ReleaseSizeOutPutMap map[string]int64) (string, string, int64){
	// Check if the user provided version is on the list we have just extracted
	if Contains(ReleaseVersion, cmdOptions.Version) {
		v := cmdOptions.Version
		Infof("Found the GPDB version \"%s\" on PivNet, continuing..", v)
		return ReleaseFileOutPutMap[v], ReleaseOutputMap[v], ReleaseSizeOutPutMap[v]
	} else { // If its not on the list then fallback to interactive mode
		// Print warning if the user did provide a value of the version
		downloadFromWhere := "PivNet"
		if cmdOptions.Github {
			downloadFromWhere = "Github"
		}
		
		if cmdOptions.Version != "" {
			Warnf("Unable to find the GPDB version \"%s\" on %s, failing back to interactive mode..",
				cmdOptions.Version, downloadFromWhere)
		}

		// Sort all the keys
		var data = []string{"Index | Product Version",
			"---------| ------------------------",
		}
		for index, version := range ReleaseVersion {
			data = append(data, strconv.Itoa(index+1)+"|"+version)
		}
		printOnScreen("Please select the version from the drop down list", data)

		// Total accepted values that user can enter
		TotalOptions := len(ReleaseVersion)

		// Ask user for choice
		usersChoice := PromptChoice(TotalOptions)

		// Selected by the user
		choice := ReleaseVersion[usersChoice-1]
		cmdOptions.Version = choice

		// Return would be the filename, url, filesize
		return ReleaseFileOutPutMap[choice], ReleaseOutputMap[choice], ReleaseSizeOutPutMap[choice]
	}
}

// Ask user what file in that version are they interested in downloading
// Default is to download GPDB, GPCC and other with a option from parser
func (r *Responses) WhichProduct(token string) {
	Info("Checking for the version to download")

	// Clearing up the buffer to ensure we are using a clean array and MAP
	var ProductOutputMap = make(map[string]string)
	var ProductOptions = []string{}

	// Crap out if the user choose GPDB version less that 5.24
	c := extractVersion(r.UserRequest.versionChoosen)
	if c > 5.24 && !r.GPCC && cmdOptions.Product == "gpcc" {
		Fatalf("GPCC is no longer connected to GPDB product category from version 5.24 onwards, retry again..")
	}

	// For command center greater than 5.23
	if cmdOptions.Product == "gpcc" && r.isThisGPCC4xAboveRequest() {
		for _, j := range r.VersionList.Product_files {
			ProductOutputMap[j.Name] = j.Links.Self.Href
			ProductOptions = append(ProductOptions, j.Name)
		}
	} else {
		// This is the correct API, all the files are inside the file group MAP
		isItUbuntu := getSystemInfoAndCheckIfItsUbuntu()
		for _, k := range r.VersionList.File_groups {
			// GPDB Options
			if cmdOptions.Product == "gpdb" {
				var rx *regexp.Regexp
				if isItUbuntu && !isThisGPDB6xAndAbove() { // if this is ubuntu release and version to download is below GP6
					Fatalf("Ubuntu greenplum download is only available with GPDB 6.x and above ")
				} else if isThisGPDB6xAndAbove() { // From version 6 we will use the newer regex
					rx, _ = regexp.Compile("(?i)" + rx_gpdb_for_6_n_above)
					if isItUbuntu { // If the os is ubuntu
						Debug("Changed the regex for finding the ubuntu product in pivnet")
						rx, _ = regexp.Compile("(?i)" + rx_gpdb_ubuntu_release)
					}
				} else {
					rx, _ = regexp.Compile("(?i)" + rx_gpdb)
				}

				for _, j := range k.Product_files {
					if rx.MatchString(j.Name) {
						Debugf("gpdb product list: %v", rx.FindString(j.Name))
						r.UserRequest.ProductFileURL = j.Links.Self.Href
						r.UserRequest.DownloadURL = j.Links.Download.Href
					}
				}
			}

			// GPCC option
			if cmdOptions.Product == "gpcc" {
				rx, _ := regexp.Compile(rx_gpcc)
				if rx.MatchString(k.Name) {
					Debugf("gpdb product list: ", rx.FindString(k.Name))
					for _, j := range k.Product_files {
						ProductOutputMap[j.Name] = j.Links.Self.Href
						ProductOptions = append(ProductOptions, j.Name)
					}
				}
			}

			// Others or fallback method
			if cmdOptions.Product == "gpextras" {
				for _, j := range k.Product_files {
					ProductOutputMap[j.Name] = j.Links.Self.Href
					ProductOptions = append(ProductOptions, j.Name)
				}
			}
		}
	}

	// If its GPCC or GPextras, then ask users for choice.
	if (cmdOptions.Product == "gpextras" || cmdOptions.Product == "gpcc") && len(ProductOptions) != 0 {
		var data = []string{"Index | Products",
			"----------| ------------------------------------------------",
		}
		for index, product := range ProductOptions {
			data = append(data, strconv.Itoa(index+1)+"|"+product)
		}

		printOnScreen("Please select the product from the drop down list", data)

		TotalOptions := len(ProductOptions)
		userChoice := PromptChoice(TotalOptions)
		versionSelectedUrl := ProductOutputMap[ProductOptions[userChoice-1]]
		r.UserRequest.ProductFileURL = versionSelectedUrl
		r.UserRequest.DownloadURL = versionSelectedUrl + "/download"
	}

	Debugf("Product File URL: %s", r.UserRequest.ProductFileURL)
	Debugf("Download File URL: %s", r.UserRequest.DownloadURL)

	// We received the download URL, lets gets the size and file name of the download file.
	r.ExtractFileNamePlusSize(token)
}