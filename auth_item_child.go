package annoying

type AuthItemChild struct {
	BaseModel
	Parent string `json:"parent"`
	Child  string `json:"child"`
}
