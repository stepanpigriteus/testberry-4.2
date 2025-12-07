package domain

type Config struct {
	After      int
	Before     int
	Context    int
	CountOnly  bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
	Filename   string
	Work       bool
	Ports      []string
	Mode       bool
}

type Line struct {
	LineNum int    `json:"lineNum"`
	Text    string `json:"text"`
}

type Conn struct {
	Port string
	Host string
}
