package annoying

import (
	"github.com/jinzhu/gorm"
	"time"
)

type BaseModel struct {
	Id        uint `gorm:"primary_key" json:"id"`
	CreatedAt uint `json:"created_at"`
	UpdatedAt uint `json:"updated_at"`
}

func (b *BaseModel) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("CreatedAt", time.Now().Unix())
	err = scope.SetColumn("UpdatedAt", time.Now().Unix())
	return
}

func (b *BaseModel) BeforeUpdate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("UpdatedAt", time.Now().Unix())
	return
}