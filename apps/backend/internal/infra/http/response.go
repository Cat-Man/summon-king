package httpx

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id"`
}

func Success(data interface{}, traceID string) APIResponse {
	return APIResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
		TraceID: traceID,
	}
}

func Error(code int, message string, traceID string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		TraceID: traceID,
	}
}
