package download

import (
  "fmt"
  "encoding/json"
  "errors"
  "strings"
  "github.com/op/go-logging"
  "github.com/ielizaga/piv-go-gpdb/core"
  "regexp"
)

var (
  log = logging.MustGetLogger("gpdb")
)

// Extract all the Pivotal Network Product from the Product API page.
func extract_product() (ProductObjects, error) {

  log.Info("Obtaining the product ID")

  // Get the API from the Pivotal Products URL
  ProductApiResponse, err := GetApi("GET", Products, false, "", 0)
  if err != nil {return ProductObjects{}, err}

  // Store all the JSON on the Product struct
  json.Unmarshal(ProductApiResponse, &ProductJsonType)

  // Return the struct
  return ProductJsonType, nil
}


// Extract all the Releases of the product with slug : pivotal-gpdb
func extract_release(productJson ProductObjects) (ReleaseObjects, error) {

  // Check what is the URL for the all the releases for product of our interest
  for _, product := range productJson.Products {
    if product.Slug == ProductName {
      ReleaseURL = product.Links.Releases.Href
      PivotalProduct = product.Name
    }
  }

  log.Info("Obtaining the releases for product: " + PivotalProduct)

  // If we do find the release URL, lets continue
  if ReleaseURL != "" {

    // Extract all the releases
    ReleaseApiResponse, err := GetApi("GET", ReleaseURL, false, "", 0)
    if err != nil {return ReleaseObjects{}, err}

    // Store all the releases on the release struct
    json.Unmarshal(ReleaseApiResponse, &ReleaseJsonType)

  } else { // Else lets error out
    return ReleaseJsonType, errors.New("Cannot find any Release URL for slug ID: " + ProductName)
  }

  // Return the release struct
  return ReleaseJsonType, nil
}

// provide choice of which version to download
func show_available_version(ReleaseJson ReleaseObjects) (string, string, error) {

  // Local storehouse
  var version_selected string
  var release string

  // Get all the releases from the ReleaseJson
  for _, release := range ReleaseJson.Release {
    ReleaseOutputMap[release.Version] = release.Links.Self.Href
    ReleaseVersion = append(ReleaseVersion, release.Version)
  }

  // Check if the user provided version is on the list we have just extracted
  if core.Contains(ReleaseVersion, core.RequestedDownloadVersion) {
    log.Info("Found the GPDB version \""+ core.RequestedDownloadVersion +"\" on PivNet, continuing..")
    version_selected = ReleaseOutputMap[core.RequestedDownloadVersion]
    release = core.RequestedDownloadVersion

  } else { // If its not on the list then fallback to interactive mode

    // Print warning if the user did provide a value of the version
    if core.RequestedDownloadVersion != "" {
      log.Warning("Unable to find the GPDB version \""+ core.RequestedDownloadVersion +"\" on PivNet, "+
                  "failing back to interactive mode..\n")
    } else { // print a blank line
      fmt.Println()
    }

    // Sort all the keys
    for index, version := range ReleaseVersion {
      fmt.Printf("%d: %s\n", index+1, version)
    }

    // Total accepted values that user can enter
    TotalOptions := len(ReleaseVersion)

    // Ask user for choice
    users_choice := core.Prompt_choice(TotalOptions)

    // Selected by the user
    version_selected = ReleaseOutputMap[ReleaseVersion[users_choice-1]]
    release = ReleaseVersion[users_choice-1]
  }

  return release, version_selected, nil
}

// From the user choice extract all the files available on that version
func extract_downloadURL(ver string, url string) (VersionObjects, error){

  log.Info("Obtaining the files under the greenplum version: " + ver)
  log.Debugf(url)

  // Extract all the files from the API
  VersionApiResponse, err := GetApi("GET", url, false, "", 0)
  core.Fatal_handler(err)

  // Load it to struct
  json.Unmarshal(VersionApiResponse, &VersionJsonType)

  // Updating the EULA Acceptance link
  EULA = VersionJsonType.Links.Eula_acceptance.Href

        // Return the result
  return VersionJsonType, nil
}

// Ask user what file in that version are they interested in downloading
// Default is to download GPDB, GPCC and other with a option from parser
func which_product(versionJson VersionObjects, VerToDownload string) error {

  // Clearing up the buffer to ensure we are using a clean array and MAP
  ProductOutputMap = make(map[string]string)
  ProductOptions = []string{}

  // This is the correct API, all the files are inside the file group MAP
  for _, k := range versionJson.File_groups {

    // GPDB Options
    if core.RequestedDownloadProduct == "gpdb" {

      rx, _ := regexp.Compile(rx_gpdb)

        for _, j := range k.Product_files {
          if (rx.MatchString(j.Name)) {
              log.Debugf(rx.FindString(j.Name))
              ProductFileURL = j.Links.Self.Href
              DownloadURL = j.Links.Download.Href
          }
        }
    }

    // GPCC option
    if core.RequestedDownloadProduct == "gpcc" {

      rx, _ := regexp.Compile(rx_gpcc)

      if (rx.MatchString(k.Name)) {
        log.Debugf(rx.FindString(k.Name))
        for _, j := range k.Product_files {
          ProductOutputMap[j.Name] = j.Links.Self.Href
          ProductOptions = append(ProductOptions, j.Name)
        }
      }
    }

    // Others or fallback method
    if core.RequestedDownloadProduct == "gpextras" {
      for _, j := range k.Product_files {
        ProductOutputMap[j.Name] = j.Links.Self.Href
        ProductOptions = append(ProductOptions, j.Name)
      }
    }

  }
  

  // If its GPCC or GPextras, then ask users for choice.
  if (core.RequestedDownloadProduct == "gpextras" || core.RequestedDownloadProduct == "gpcc") && len(ProductOptions) != 0 {
    fmt.Println()
    for index, product := range ProductOptions {
      fmt.Printf("%d: %s\n", index+1, product)
    }
    TotalOptions := len(ProductOptions)
    users_choice := core.Prompt_choice(TotalOptions)
    version_selected_url := ProductOutputMap[ProductOptions[users_choice-1]]
    ProductFileURL = version_selected_url
    DownloadURL = version_selected_url + "/download"
  }
  return nil
}

// extract the filename and the size of the product that the user has choosen
func extract_filename_and_size (url string) error {

  log.Info("Extracting the filename and the size of the product file.")

  // Obtain the JSON response of the product file API
  ProductFileApiResponse, err := GetApi("GET", url, false , "" , 0)
  if err != nil {return err}

  // Store it on JSON
  json.Unmarshal(ProductFileApiResponse, &ProductFileJsonType)

  // Get the filename and the file size
  filename := strings.Split(ProductFileJsonType.Product_file.Aws_object_key, "/")
  ProductFileName = filename[len(filename)-1]
  ProductFileSize = ProductFileJsonType.Product_file.Size

  log.Info("filename:" + ProductFileJsonType.Product_file.Aws_object_key)
        log.Info("ProductFileName:" + ProductFileName)

  return err

}

func Download() error {

  // Authentication validations
  log.Info("Checking if the user is a valid user")
  _, err := GetApi("GET", Authentication, false, "", 0)
  if err != nil {return err}

  // Product list
  productJson, err := extract_product()
  if err != nil {return err}

  // Release list
  releaseJson, err := extract_release(productJson)
  if err != nil {return err}

  // What is user's choice
  choice, choice_url, err := show_available_version(releaseJson)
  if err != nil {return err}

  // Set the of Install if the user has asked to install after download and not version information is available
  if core.InstallAfterDownload { core.RequestedInstallVersion = choice }

  // Get all the files under that version
  versionFileJson, err := extract_downloadURL(choice, choice_url)
  if err != nil {return err}

  // The users choice to what to download from that version
  err = which_product(versionFileJson, choice)
  if err != nil {return err}

  // If we didn't find the database File, then fall back to interactive mode.
  if (core.RequestedDownloadProduct == "gpdb" || core.RequestedDownloadProduct == "gpcc") && ProductFileURL == "" {
    log.Warning("Couldn't find binaries for GPDB version \"" + choice + "\", failing back to interactive mode...")
    core.RequestedDownloadProduct = "gpextras"
    which_product(versionFileJson, choice)
  }

  // Extract the filename and the size of the file
  err = extract_filename_and_size(ProductFileURL)
  if err != nil {return err}

  // Accept the EULA
  log.Info("Accepting the EULA (End User License Agreement): " + EULA)
  _, err = GetApi("POST", EULA, false, "", 0)
  if err != nil {return err}

  // Download the version
  log.Info("Starting downloading of file: " + ProductFileName)
  if DownloadURL != "" {
    _, err := GetApi("GET", DownloadURL, true, ProductFileName, ProductFileSize)
    if err != nil {return err}
  } else {
    return errors.New("Download URL is blank, cannot download the product.")
  }

  return nil
}
