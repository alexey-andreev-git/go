package entity

// EntityReference represents the reference data associated with an entity in the database
type EntityReference struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}
type EntityReferenceById struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// EntityDataReference represents the reference data associated with an entity in the database
type EntityDataReference struct {
	Reference int    `json:"reference"`
	Order     int    `json:"order"`
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
	Id        int `json:"id"`
	Reference int `json:"reference"`
}

// EntityData represents the data associated with an entity in the database
type EntityData struct {
	Entity int `json:"entity"`
	Order  int `json:"order"`
	Value  int `json:"value"`
}
