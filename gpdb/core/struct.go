package core

import (
	"github.com/codingconcepts/env"
)

type Environment struct {
	BaseDir string `env:"BASE_DIR"`
	TempDir string `env:"TEMP_DIR"`
	AppName string `env:"APPLICATION_NAME"`
}

type Download struct {
	ApiToken    string `env:"API_TOKEN"`
	DownloadDir string `env:"DOWNLOAD_DIR"`
}

type Install struct {
	EnvDir               string `env:"ENV_DIR"`
	UnistallDir          string `env:"UNINTSALL_DIR"`
	FutureRef            string `env:"FUTUREREF_DIR"`
	MasterHost           string `env:"MASTER_HOST"`
	MasterUser           string `env:"MASTER_USER"`
	MasterPass           string `env:"MASTER_PASS"`
	GpMonPass            string `env:"GPMON_PASS"`
	MasterDataDirectory  string `env:"MASTER_DATA_DIRECTORY"`
	SegmentDataDirectory string `env:"SEGMENT_DATA_DIRECTORY"`
	TotalSegments        int    `env:"TOTAL_SEGMENT"`
}
