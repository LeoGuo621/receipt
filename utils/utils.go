package utils

type Resp struct {
	Recode string `json:"recode"`
	Msg string `json:"msg"`
	// valid data
	Data interface{} `json:"data"`
}

// scheduled some predictable err recodes
const (
	RecodeOk        = "0"
	RecodeDBErr     = "4001"
	RecodeLoginErr  = "4002"
	RecodeParamErr  = "4003"
	RecodeSysErr    = "4004"
	RecodeEthErr    = "4105"
	RecodeUnknownErr = "4106"
)

var recodeText = map[string]string{
	RecodeOk:        "Success",
	RecodeDBErr:     "Database operation error",
	RecodeLoginErr:  "User login error",
	RecodeParamErr:  "Parameter error",
	RecodeSysErr:    "System error",
	RecodeEthErr:   "Error interacting with Ethereum",
	RecodeUnknownErr: "Unknown error",
}

func RecodeText(code string) string {
	if str, ok := recodeText[code]; ok {
		return str
	}
	return recodeText[RecodeUnknownErr]
}
