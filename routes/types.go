package routes

/*
JSon struct to be sent to jquery file upload widget
*/

// collaction of the file informations uploaded and traited
// swagger:response  fileInfos
type FileInfos struct {
	// list of the files analyzed
	// swagger:allOf
	Files []*FileInfo `json:"files"`
}

// informations about the uploaded file
// swagger:response fileInfo
type FileInfo struct {
	// name of the file
	Name string `json:"name"`
	// size of the file
	Size int64 `json:"size"`
	// error string encountered
	Error string `json:"error,omitempty"`
}
