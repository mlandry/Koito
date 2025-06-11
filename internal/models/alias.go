package models

type Alias struct {
	ID      int32  `json:"id"`
	Alias   string `json:"alias"`
	Source  string `json:"source"`
	Primary bool   `json:"is_primary"`
}
