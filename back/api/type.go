package api

type BookEntry struct {
	Type            string  `json:"type"`
	Path            string  `json:"path"`
	Cover           string  `json:"cover"`
	Title           string  `json:"title"`
	CurrentPosition string  `json:"currentPosition"`
	Progress        float64 `json:"progress"`
}
