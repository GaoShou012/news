package im

type Head struct {
	Id           uint64 `json:"id"`
	BusinessType string `json:"businessType"`
	BusinessApi  string `json:"businessApi"`
}
type Message struct {
	Head `json:"head"`
	Body interface{} `json:"body"`
}
