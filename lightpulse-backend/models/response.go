package models

// Response 共通レスポンスモデル
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse エラーレスポンス
func ErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Error:   message,
	}
}

// SuccessResponse 成功レスポンス
func SuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

