package install


// Install GPDB
const InstallGPDBBashFileName = "install_software.sh"
const DBName = "gpadmin"

// Port Base store house
const PortBaseFileName = "gpdbportbase.txt"
const PORT_BASE = 30000
const MASTER_PORT = 3000

// gpinitsystem temp file
const GpInitSystemConfigFile = "gpinitsystemconfig"
const MachineFileList = "hostfile"

// GPCC Port Base
const GPCCPortBaseFileName = "gpccportbase.txt"
const GPCC_PORT_BASE = 28080

// SingleNode GPDB Install
var MasterDIR string
var SegmentDIR string
var GpInitSystemConfig GpInitSystemConfigObject
var MasterHostFileList string
var BinaryInstallLocation string
var GpInitSystemConfigDir string

// Core
var EnvFileName string
var ThisDBMasterPort int
var CoreIP string = "192.0.0.0"

// GPCC
var GPCC_PORT string
var InstallWLM bool = false
var GPCC_INSTANCE_NAME string
var ThisEnvGPCCPort int
var GPPERFMONHOME string
var WLMVersion string
var WLMInstallDir string
