package core


type EnvObjects struct {
	Core coreType `yaml:"CORE"`
	Download downloadType `yaml:"DOWNLOAD"`
	Install installType `yaml:"INSTALL"`
}

type coreType struct {
	BaseDir string `yaml:"BASE_DIR"`
	TempDir string `yaml:"TEMP_DIR"`
	AppName string `yaml:"APPLICATION_NAME"`
}

type downloadType struct {
	ApiToken string `yaml:"API_TOKEN"`
	DownloadDir string `yaml:"DOWNLOAD_DIR"`
}

type installType struct {
	EnvDir string `yaml:"ENV_DIR"`
	UnistallDir string `yaml:"UNINTSALL_DIR"`
	FutureRef string `yaml:"FUTUREREF_DIR"`
	MasterHost string `yaml:"MASTER_HOST"`
	MasterUser string `yaml:"MASTER_USER"`
	MasterPass string `yaml:"MASTER_PASS"`
	GpMonPass string `yaml:"GPMON_PASS"`
	MasterDataDirectory string `yaml:"MASTER_DATA_DIRECTORY"`
	SegmentDataDirectory string `yaml:"SEGMENT_DATA_DIRECTORY"`
	TotalSegments int `yaml:"TOTAL_SEGMENT"`
}

