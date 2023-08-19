package backend

type Update struct {
	NewUpdate   bool   `json:"new_update"`
	Version     string `json:"version"`
	Updating    bool   `json:"updating"`
	Downloading int    `json:"downloading"`
	Installing  int    `json:"installing"`
}
