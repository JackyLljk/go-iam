package mysql

import (
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
	"j-iam/internal/apiserver/store"
	genericoptions "j-iam/internal/pkg/options"
)

type datastore struct {
	db *gorm.DB

	// 如果需要，可以包括两个 MySQL 实例
	//docker *gorm.DB
	//db *gorm.DB
}

func (d *datastore) Users() store.UserStore { return newUsers(d) }

func (d *datastore) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

// GetMysqlFactory 使用给定的配置创建 MySQL 工厂
func GetMysqlFactory(opts *genericoptions.InsecureServingOptions)
