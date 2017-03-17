package routes

/*
JSon struct to be sent to jquery file upload widget
*/

// swagger:response
type FileInfos struct {
	// list of the files analyzed
	Files []*FileInfo `json:"files"`
}

// swagger:response
type FileInfo struct {
	// name of the file
	Name string `json:"name"`
	// size of the file
	Size int64 `json:"size"`
	// error string encountered
	Error string `json:"error,omitempty"`
}
