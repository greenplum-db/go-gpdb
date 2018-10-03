package download

// All the PivNet Url's
const (
	EndPoint = "https://network.pivotal.io"
	Authentication = EndPoint + "/api/v2/authentication"
	RefreshToken = EndPoint + "/api/v2/authentication/access_tokens"
	Products =  EndPoint + "/api/v2/products"
)

// Slug ID that is of our importance
const ProductName = "pivotal-gpdb"

// Common Downloading options
const (
	DBServer = "Database Server"
	GPCC = "Greenplum Command Center"
)

// Products
var ProductJsonType ProductObjects

// Authentication
var AuthenicationResponse AuthResp

// Release
var ReleaseJsonType ReleaseObjects
var ReleaseURL string
var PivotalProduct string
var ReleaseOutputMap = make(map[string]string)
var ReleaseVersion []string

// All files of the selected version
var VersionJsonType VersionObjects
var DowloadOutputMap = make(map[string]string)
var DownloadOption []string
var DownloadURL string
var ProductFileURL string
var ChoiceMap VersionObjects

// Product file
var ProductFileJsonType ProductFilesObjects
var ProductFileName string
var ProductFileSize int64
var EULA string
var ProductOutputMap = make(map[string]string)
var ProductOptions []string
var FileNameContains = []string{
	"Red Hat Enterprise Linux",
	"RedHat Entrerprise Linux",
	"RedHat Enterprise Linux",
	"REDHAT ENTERPRISE LINUX",
	"RHEL",
	"Binary Installer for RHEL"}

var IgnoreFileExtension = "RPM for RHEL"