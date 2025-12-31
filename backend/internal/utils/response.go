package utils

import (
        "encoding/json"
        "net/http"
)

type APIResponse struct {
        Success bool        `json:"success"`
        Data    interface{} `json:"data,omitempty"`
        Error   *APIError   `json:"error,omitempty"`
        Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
        Code    string            `json:"code"`
        Message string            `json:"message"`
        Details map[string]string `json:"details,omitempty"`
}

type Meta struct {
        RequestID string `json:"request_id,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(APIResponse{
                Success: status >= 200 && status < 300,
                Data:    data,
        })
}

func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, requestID string) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(APIResponse{
                Success: status >= 200 && status < 300,
                Data:    data,
                Meta:    &Meta{RequestID: requestID},
        })
}

func ErrorResponse(w http.ResponseWriter, status int, code, message string, details map[string]string) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(APIResponse{
                Success: false,
                Error: &APIError{
                        Code:    code,
                        Message: message,
                        Details: details,
                },
        })
}

func BadRequest(w http.ResponseWriter, message string, details map[string]string) {
        ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(w http.ResponseWriter, message string) {
        ErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(w http.ResponseWriter, message string) {
        ErrorResponse(w, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(w http.ResponseWriter, message string) {
        ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func InternalError(w http.ResponseWriter, message string) {
        ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

func Conflict(w http.ResponseWriter, message string) {
        ErrorResponse(w, http.StatusConflict, "CONFLICT", message, nil)
}
