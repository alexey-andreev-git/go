package models

type (
	ControlData struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		Placeholder string   `json:"placeholder"`
		Validators  []string `json:"validators"`
	}
	Controls struct {
		Controls []ControlData `json:"controls"`
	}
)
