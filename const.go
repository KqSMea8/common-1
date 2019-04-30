package common

//一些常量数据
const SERVICE_NAME = "zmind.api"
const SERVICE_NAME_MORPHEUS = "morpheus.api.FileAttribute"
const SERVICE_NAME_MORPHEUS_STORAGE = "morpheus.api.Storage"
const SERVICE_NAME_MORPHEUS_DETECTION = "morpheus.api.FileDetection"
const SERVICE_NAME_KAMALA = "kamala.api.job"
const SERVICE_NAME_INTERAL_RPC = "127.0.0.1:7890"
const TIME_FORMAT = "2006-01-02 15:04:05"

const DOWN_LOAD_PATH = "/opt/zmind-service/cuckoo/storage/analyses"

var FILE_TYPE_ASSET = map[string]string{
	"typer": "typer",
	"qex":   "qex",
}
var EXT_DESC = map[string]string{
	"nsis_exe": "NSIS Script",
	"nsisexe":  "NSIS Script",
	"exe32":    "PE32 executable for MS Windows (GUI) Intel 80386 32-bit",
	"exe64":    "PE64 executable for MS Windows (GUI)",
	"peexe":    "PE executable for MS Windows (GUI)",
	"dll32":    "PE32 executable for MS Windows (DLL) (console) Intel 80386 32-bit",
	"dll64":    "PE64 executable for MS Windows (DLL)",
	"pedll":    "PE executable for MS Windows (DLL)",
	"pdf":      "PDF document",
	"jar":      "Java Archive",
	"xls":      "Microsoft Excel sheet",
	"doc":      "Microsoft Word document",
	"ppt":      "Microsoft PowerPoint document",
	"pps":      "PowerPoint Slide Show",
	"xlw":      "Excel Workspace File",
	"xlsx":     "Excel Microsoft Office Open XML Format document",
	"xlsm":     "Excel Open XML Macro-Enabled Presentation",
	"xlsb":     "Excel Binary Spreadsheet",
	"xlam":     "Excel Open XML Macro-Enabled Add-In",
	"docx":     "Word Microsoft Office Open XML Format document",
	"docm":     "Word Open XML Macro-Enabled Presentation",
	"dotm":     "Word Open XML Macro-Enabled Document Template",
	"dotx":     "Word Open XML Document Template",
	"pot":      "PowerPoint Template",
	"pptx":     "PowerPoint Microsoft Office Open XML Format document",
	"ppsx":     "PowerPoint Open XML Slide Show",
	"ppam":     "PowerPoint 2007 Add-In",
	"pptm":     "PowerPoint Open XML Macro-Enabled Presentation",
	"ppsm":     "PowerPoint Open XML Macro-Enabled Slide Show",
	"html":     "HyperText Markup Language",
	"zip":      "ZIP compressed archive",
	"rar":      "RAR archive data",
	"7z":       "7-zip archive data",
	"lnk":      "MS Windows shortcut",
	"swf":      "Macromedia Flash Player Compressed Movie",
	"bat":      "DOS batch file text",
	"mht":      "MIME entity text",
	"eml":      "RFC 822 mail text",
	"msi":      "Microsoft Windows Installer",
	"wsf":      "Windows Script File",
	"msg":      "Outlook Message",
	"js":       "JavaScript File",
	"vbs":      "VBScript File",
	"ps1":      "Windows PowerShell Cmdlet File",
	"rtf":      "Rich Text Format File",
	"py":       "Python script text executable",
	"gz":       "gzip compressed data",
}

var VUL_FILE_EXTS = []string{"exe", "dll", "sys",
	"rtf", "ole", "pdf", "doc", "docm", "docx", "xls", "xlsm", "xlsx", "ppt",
	"pptm", "pptx", "eml", "mht", "swf", "jar", "wsf", "msi", "nsis_exe", "mz", "nsis",
	"7z", "rar", "zip", "bat", "ps1", "vbs", "py", "js", "html", "xml", "lnk", "htm/html", "nsisexe"}

//每种签订器可以扫描的文件类型
var FILE_TYPE_TO_RIGHT_ASSET = map[string][]string{
	"bole":   {"exe", "nsisexe", "dll", "sys"},
	"sign":   {"exe", "nsisexe", "dll", "sys"},
	"owl":    {"exe", "nsisexe", "dll", "sys", "lnk", "ole", "rtf", "chm", "vbe", "pdf", "swf", "zip", "active_mime", "mail", "elf", "macho"},
	"cuckoo": VUL_FILE_EXTS,
}

// qex 日志结构
type QexInfo struct {
	EmbeddedLayer    int32  `json:"embedded_layer"`
	FilePath         string `json:"file_path"`
	OriginalFileName string `json:"original_file_name"`
	EmbeddeFileName  string `json:"embedded_file_name"`
	VirusMajorType   int32  `json:"virus_major_type"`
	VirusMinorType   int32  `json:"virus_minor_type"`
	FileType         string `json:"file_type"`
	FirstAction      int32  `json:"first_action"`
	SecondAction     int32  `json:"second_action"`
	ErrorCOde        int32  `json:"error_code"`
	FileTypeCallback string `json:"errofile_type_callbackr_code"`
	MalwareName      string `json:"virus_name"`
	Status           int    `json:"scan_status"`
}

type CERT struct {
	Organization string `json:"organization"`

	EndTime      string `json:"end_time"`
	Country      string `json:"country"`
	CommonName   string `json:"common_name"`
	SerialNumber string `json:"serial_number"`
	StartTime    string `json:"start_time"`
	Locality     string `json:"locality"`
	Email        string `json:"email"`
	Md5          string `json:"md5"`
	Sha1         string `json:"sha1"`
}
type HASH struct {
	Md5    string `json:"md5"`
	Sha256 string `json:"sha256"`
	Ssdeep string `json:"ssdeep"`
	Sha512 string `json:"sha512"`
	Sha1   string `json:"sha1"`
}

type ICONINFO struct {
	Png string `json:"png"`
}

// typer日志结构

type TyperInfo struct {
	Exiftool []string `json:"exiftool"`
	Trid     []string `json:"trid"`
	Cert     []CERT   `json:"cert"`
	Hash     HASH     `json:"hash"`
	Icon     ICONINFO `json:"icon"`
}

// sign日志结构
type SignInfo struct {
	MalwareName string `json:"cn"`
	Status      int    `json:"sign_status"`
	Sign        string `json:"sign_hash"`
}

//qvm日志结构
type QvmInfo struct {
	MalwareName string `json:"virus_name"`
	Status      int    `json:"judge"`
	Model       int
	Score       float32
}

//ave日志结构
type AveInfo struct {
	MalwareName string `json:"virus_name"`
	Status      int    `json:"scan_reslut"`
}

// db日志结构
type BdInfo struct {
	MalwareName string `json:"threat_info"`
	Status      int
	ThreatType  int
}

// owl引擎结构
type OwlInfo struct {
	FilePath string       `json:"filepath"`
	OwlTags  []string     `json:"owl_tags"`
	Streams  []StreamInfo `json:"stream_info"`
}
type StreamInfo struct {
	Md5       string
	Sha1      string
	Pdb       string
	Bit       int            `json:"64bits"`
	Name      string         `json:"stream_name"`
	Stype     string         `json:"pe_sub_type"`
	FType     string         `json:"stream_type"`
	CType     string         `json:"compiler_type"`
	Header    HeaderInfo     `json:"pe_header_info"`
	Export    ExportInfo     `json:"exports_info"`
	Imports   []ImportInfo   `json:"imports_info"`
	Resources []ResourceInfo `json:"sres_info"`
	Versions  []VersionInfo  `json:"version_info"`
}
type HeaderInfo struct {
	Time     int64         `json:"timestamp_int"`
	Sections []SectionInfo `json:"section"`
}
type SectionInfo struct {
	Name string
}
type VersionInfo struct {
	LC string `json:"legalcopyright"`
	IN string `json:"internalname"`
	FV string `json:"fileversion"`
	LG string `json:"Language"`
	CN string `json:"companyname"`
	PN string `json:"productname"`
	PV string `json:"productversion"`
	FD string `json:"filedescription"`
	ON string `json:"originalfilename"`
}
type ExportInfo struct {
	Name string   `json:"exportname"`
	Apis []string `json:"exports_api"`
}
type ImportInfo struct {
	Name string   `json:"libname"`
	Apis []string `json:"apiname"`
}
type ResourceInfo struct {
	SType string   `json:"rt_type"`
	Infos []ReInfo `json:"info"`
}
type ReInfo struct {
	Name string
	Size int
}

//文件基本属性结构
type FileBaseInfo struct {
	Name string
	Size uint32
	Sha1 string
}

type CuckooInfo struct {
	Pcap     string
	Snapshot []string
}
