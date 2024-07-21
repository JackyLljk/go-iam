package mysql

import (
	"context"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"gorm.io/gorm"
	v1 "j-iam/pkg/model/v1"
)

type users struct {
	db *gorm.DB
}

func newUsers(d *datastore) *users { return &users{d.db} }

// Create 创建一个新的用户账号
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.db.Create(&user).Error
}
