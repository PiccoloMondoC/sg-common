// sg-common/logtypes/logtypes.go
package logtypes

type LogEntry struct {
	Service     string                 `json:"service"`
	LogName     string                 `json:"logName"`
	Resource    Resource               `json:"resource"`
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	InsertId    string                 `json:"insertId,omitempty"`
	HttpRequest HttpRequest            `json:"httpRequest,omitempty"`
	TextPayload string                 `json:"textPayload,omitempty"`
	JsonPayload map[string]interface{} `json:"jsonPayload,omitempty"`
}

type Resource struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type HttpRequest struct {
	RequestMethod string `json:"requestMethod"`
	RequestUrl    string `json:"requestUrl"`
	Status        int    `json:"status"`
	UserAgent     string `json:"userAgent"`
	RemoteIp      string `json:"remoteIp"`
}

type LogData struct {
	Log LogEntry `json:"log"`
}
