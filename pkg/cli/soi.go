package cli

type SoiData struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	URI     string   `json:"uri"`
	Tags    []string `json:"tags"`
	Created string   `json:"created"`
}

type SoiWithPath struct {
	*SoiData
	Path string `json:"path"`
}

type SoiBucket struct {
	Sois []*SoiWithPath `json:"sois"`
}

// TODO implement later
type Client struct {
	WorkingDir string
}
