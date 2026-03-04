package responseModel

type Response struct {
	Type         string      `json:"type"`
	Message      string      `json:"message,omitempty"`
	Error        interface{} `json:"error,omitempty"`
	PageMetaData interface{} `json:"pageMetaData,omitempty"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
	*Response
}
