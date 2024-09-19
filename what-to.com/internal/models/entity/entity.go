package entity

// EntityReference represents the reference data associated with an entity in the database
type EntityReference struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}
type EntityReferenceById struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// EntityDataReference represents the reference data associated with an entity in the database
type EntityDataReference struct {
	Reference uint   `json:"reference"`
	Order     uint   `json:"order"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Comment   string `json:"comment"`
}
type EntityDataReferenceByOrder struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

// Entity represents an entity in the database
type Entity struct {
	Id        uint `json:"id"`
	Reference uint `json:"reference"`
}

// EntityData represents the data associated with an entity in the database
type EntityData struct {
	Entity uint `json:"entity"`
	Order  uint `json:"order"`
	Value  uint `json:"value"`
}
