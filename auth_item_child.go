package annoying

type AuthItemChild struct {
	BaseModel
	Parent string `json:"parent"`
	Child  string `json:"child"`
}

func (*AuthItemChild) TableName() string {
	return "auth_item_child"
}
