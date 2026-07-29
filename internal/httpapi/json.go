package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func DecodeJSONFromBytes(b []byte, v any) error {
	dec := json.NewDecoder(io.NopCloser(bytesReader{b: b}))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

type bytesReader struct{ b []byte }

func (r bytesReader) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
