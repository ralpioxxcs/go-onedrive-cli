package graph

type BaseItem struct {
	Id                   string      `json:"id"`
	CreatedBy            string      `json:"createdBy"`
	CreatedDateTime      string      `json:"createdDataTime"`
	Etag                 string      `json:"eTag"`
	LastModifiedBy       IdentitySet `json:"lastModifiedBy"`
	LastModifiedDateTime string      `json:"lastModifiedDateTime"`
	Name                 string      `json:"name"`
	//ParentReference
	//WebURL
}

// IdentitySet represents a key collection of identity resources
type IdentitySet struct {
	Application Identity `json:"application"`
	Device      Identity `json:"device"`
	User        Identity `json:"user"`
}

type Identity struct {
	DisplayName string `json:"displayName"` // display name of identity
	Id          string `json:"id"`          // unique identifier for the identity
}
