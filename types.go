package main

type Configs struct {
	SrcDB DBConfigs  `json:"srcdb"`
	OutDB DBConfigs  `json:"outdb"`
	List  []Contents `json:"list"`
}

type DBConfigs struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type Contents struct {
	Mode   string   `json:"mode"`
	Hash   string   `json:"hash"`
	Key    string   `json:"key"`
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Bypass []string `json:"bypass"`
}
