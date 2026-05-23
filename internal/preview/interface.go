package preview

import "net/http"

type PreviewService interface {
	ServePreview(w http.ResponseWriter, r *http.Request)
}
