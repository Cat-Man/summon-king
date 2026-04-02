package http

// APIResponse wraps every API response with traceability.
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id"`
}

// Success returns a standard success payload.
func Success(data interface{}, traceID string) APIResponse {
    return APIResponse{Code: 0, Message: "ok", Data: data, TraceID: traceID}
}

// Error returns a structured error payload.
func Error(code int, message, traceID string) APIResponse {
    return APIResponse{Code: code, Message: message, TraceID: traceID}
}
