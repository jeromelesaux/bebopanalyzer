package routes

/*
JSon struct to be sent to jquery file upload widget
*/
type FileInfos struct {
	Files []*FileInfo `json:"files"`
}

type FileInfo struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Error string `json:"error,omitempty"`
}
