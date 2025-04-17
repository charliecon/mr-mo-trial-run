package mock_dynamo

type Table struct {
	Items []Item `json:"items"`
}

type Item struct {
	ResourceType   string       `json:"resourceType"`
	SourceEntityId string       `json:"sourceEntityId"`
	TargetInfo     []TargetInfo `json:"targetInfo"`
}

type TargetInfo struct {
	OrgId          string `json:"orgId"`
	TargetEntityId string `json:"targetEntityId"`
}
