package annoying

type AuthAssignment struct {
	BaseModel
	ItemName string `json:"item_name"`
	UserId   string `json:"user_id"`
}

func (*AuthAssignment) TableName() string {
	return "auth_assignment"
}