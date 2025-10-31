package api

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
)

// ErrorResponse is the structure for our standard JSON error.
type ErrorResponse struct {
    Error struct {
        Message string                 `json:"message"`
        Details string                 `json:"details,omitempty"`
        Code    int                    `json:"code"`
        ID      string                 `json:"id,omitempty"`
        Meta    map[string]interface{} `json:"meta,omitempty"`
    } `json:"error"`
}

// JSONError sends a structured JSON error response.
// showDetails should be true to include backend details in the response.
func JSONError(w http.ResponseWriter, message string, details string, code int, meta map[string]interface{}, showDetails bool) {
    resp := ErrorResponse{}
    resp.Error.Message = message
    if showDetails {
        resp.Error.Details = details
    }
    resp.Error.Code = code
    resp.Error.ID = uuid.New().String()
    if meta != nil {
        resp.Error.Meta = meta
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        // fallback to plain text
        http.Error(w, message, code)
    }
}
