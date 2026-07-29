package payments

import (
	"encoding/json"
	"net/http/httptest"
)

func decodeJSONFromRecorder(w *httptest.ResponseRecorder, dst any) error {
	return json.Unmarshal(w.Body.Bytes(), dst)
}
