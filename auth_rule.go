package annoying

type AuthRule struct {
	BaseModel
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (AuthRule) TableName() string {
	return "auth_rule"
}
