package download

type LinksType struct {
	Self   HrefType `json:"self"`
	Releases   HrefType `json:"releases"`
	Product_files   HrefType `json:"product_files"`
	File_groups   HrefType `json:"file_groups"`
	Signature_file_download HrefType `json:"signature_file_download"`
	Eula_acceptance HrefType `json:"eula_acceptance"`
	User_groups HrefType `json:"user_groups"`
	Download HrefType `json:"download"`
}

type HrefType struct {
	Href string `json:"href"`
}

type eulaType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}

type ProductFilesObjects struct {
	Product_file ProductFilesObjType `json:"product_file"`
}

type ProductFilesObjType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	Description string `json:"description"`
	Docs_url string `json:"docs_url"`
	File_transfer_status string `json:"file_transfer_status"`
	File_type string `json:"file_version"`
	Has_signature_file string `json:"has_signature_file"`
	Included_files []string `json:"included_files"`
	Md5 string `json:"md5"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Ready_to_serve bool `json:"ready_to_serve"`
	Released_at string `json:"released_at"`
	Size int64 `json:"size"`
	System_requirements []string `json:"system_requirements"`
	Links LinksType `json:"_links"`
}

type ProductObjects struct {
	Products []ProductObjType `json:"products"`
	Links LinksType `json:"_links"`
}

type ProductObjType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Logo_url string `json:"logo_url"`
	Links   LinksType `json:"_links"`
}

type ReleaseObjects struct {
	Release []releaseObjType `json:"releases"`
	Links LinksType `json:"_links"`
}

type releaseObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Release_notes_url string `json:"release_notes_url"`
	Availability string `json:"availability"`
	Description string `json:"description"`
	Eula   eulaType `json:"eula"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}

type VersionObjects struct {
	VersionObjType
}

type VersionObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Availability string `json:"availability"`
	Eula   eulaType `json:"eula"`
	End_of_support_date string `json:"end_of_support_date"`
	End_of_guidance_date string `json:"end_of_guidance_date"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Product_files []verProdType `json:"product_files"`
	File_groups []verFileGroupType `json:"file_groups"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}

type verProdType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	File_version string `json:"file_version"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}

type verFileGroupType struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Product_files []verProdType `json:"product_files"`
}

type AuthResp struct {
	Token string `json:"access_token"`
}

type AuthBody struct {
	RefreshToken string `json:"refresh_token"`
}
