package model

type Group struct {
	GroupKey  string   `json:"groupKey"`
	GroupName string   `json:"groupName"`
	Members   []Member `json:"members"`
}

type Member struct {
	Key         string      `json:"key"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"` //"value": [1593674764.889,"1"]
}
