package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
)

const (
	AlertMsgGeneric = "Something went wrong, please try again"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Yield interface{}
}
