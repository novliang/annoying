package annoying

const TypeRole = 1;

const TypePermission = 2;

type AuthItem struct {
	BaseModel
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	RuleName    string `json:"rule_name"`
	Data        []byte `json:"data"`
}

func (a *AuthItem) GetName() string {
	return a.Name
}

func (AuthItem) TableName() string {
	return "auth_item"
}
