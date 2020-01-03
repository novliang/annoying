package annoying

import (
	"github.com/jinzhu/gorm"
	"time"
)

const TypeRole = 1;

const TypePermission = 2;

type AuthItem struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	RuleName    string `json:"rule_name"`
	Data        []byte `json:"data"`
	CreatedAt   uint   `json:"created_at"`
	UpdatedAt   uint   `json:"updated_at"`
}

func (i *AuthItem) GetName() string {
	return i.Name
}

func (*AuthItem) TableName() string {
	return "auth_item"
}

func (i *AuthItem) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("CreatedAt", time.Now().Unix())
	err = scope.SetColumn("UpdatedAt", time.Now().Unix())
	return
}

func (i *AuthItem) BeforeUpdate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UpdatedAt", time.Now().Unix())
	return
}
