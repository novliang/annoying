package annoying

import (
	"github.com/jinzhu/gorm"
	"time"
)

type AuthAssignment struct {
	ItemName string `json:"item_name"`
	UserId   string `json:"user_id"`
	CreatedAt uint `json:"created_at"`
}

func (*AuthAssignment) TableName() string {
	return "auth_assignment"
}

func (b *AuthAssignment) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("CreatedAt", time.Now().Unix())
	return
}
