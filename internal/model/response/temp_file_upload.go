package responseModel

// TempFileUpload is the response for a single temp file upload.
// For bulk uploads: response array order matches request order (index i = i-th file).
// Use fileName to match when filenames are unique.
type TempFileUpload struct {
	Id          uint   `json:"id,omitempty"`
	TempFileKey string `json:"tempFileKey,omitempty"`
	Kind        string `json:"kind,omitempty"`
	FileName    string `json:"fileName,omitempty"`
}
