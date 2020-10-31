package soi2

type SoiData struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	URI     string   `json:"uri"`
	Tags    []string `json:"tags"`
	Created string   `json:"created"`
}
