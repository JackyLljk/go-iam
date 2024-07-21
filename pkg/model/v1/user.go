package v1

import (
	"fmt"
	"github.com/marmotedu/component-base/pkg/auth"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"gorm.io/gorm"
	"time"
)

// User 表示 user RESTful 资源，同时也用于 gorm
type User struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status int `json:"status" gorm:"column:status" validate:"omitempty"`

	// Required: true
	Nickname string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" gorm:"column:password" validate:"required"`

	// Required: true
	Email string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" gorm:"column:phone" validate:"omitempty"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`

	LoginedAt time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

// UserList 是已存储在 Storage 中的所有用户的完整列表
type UserList struct {
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

func (u *User) TableName() string { return "user" }

// Compare 与纯文本密码进行比较。如果它与加密的相同（在“User”结构中）则返回 true
func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}
	return nil
}

// AfterCreate 创建数据库记录后运行
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}
