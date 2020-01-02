package annoying

type AuthRule struct {
	BaseModel
	Name string `json:"name"`
	Data []byte `json:"data"`
}
