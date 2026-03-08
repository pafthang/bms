package encoding

import (
	"encoding/json"
	"net/http"

	"github.com/pafthang/arc"
)

type ProblemJSONEncoder struct{}

type problemDoc struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail,omitempty"`
	Code      string `json:"code,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

func (ProblemJSONEncoder) Encode(w http.ResponseWriter, status int, body any) error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/problem+json")
	}
	w.WriteHeader(status)

	apiErr, ok := body.(*arc.APIError)
	if !ok {
		return json.NewEncoder(w).Encode(body)
	}
	title := http.StatusText(status)
	if title == "" {
		title = "Error"
	}
	doc := problemDoc{
		Type:      "about:blank",
		Title:     title,
		Status:    status,
		Detail:    apiErr.Message,
		Code:      apiErr.Code,
		RequestID: apiErr.RequestID,
	}
	return json.NewEncoder(w).Encode(doc)
}
