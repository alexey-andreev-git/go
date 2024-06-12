package entity

// Entity represents an entity in the database
type Entity struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// EntityData represents the data associated with an entity in the database
type EntityData struct {
	Entity int `json:"entity"`
	Order  int `json:"order"`
	Value  int `json:"value"`
}

// EntityDataReference represents the reference data associated with an entity in the database
type EntityDataReference struct {
	Entity  int    `json:"entity"`
	Order   int    `json:"order"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}
