package mysql

import (
	"fmt"
	"j-iam/internal/apiserver/store"
	"j-iam/internal/pkg/logger"
	genericoptions "j-iam/internal/pkg/options"
	"j-iam/pkg/db"
	"sync"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB

	// 如果需要，可以包括两个 MySQL 实例
	//docker *gorm.DB
	//db *gorm.DB
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func (d *datastore) Users() store.UserStore {
	return newUsers(d)
}

func (d *datastore) Secrets() store.SecretStore {
	return newSecrets(d)
}

func (d *datastore) Policies() store.PolicyStore {
	return newPolicies(d)
}

func (d *datastore) PolicyAudits() store.PolicyAuditStore {
	return newPolicyAudits(d)
}

func (d *datastore) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

// GetMySQLFactory 使用给定的配置创建 MySQL 工厂
func GetMySQLFactory(opts *genericoptions.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			Logger:                logger.New(opts.LogLevel),
		}
		dbIns, err = db.New(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

// cleanDatabase tear downs the database tables.
//
//nolint:unused // may be reused in the feature, or just show a migrate usage.
func cleanDatabase(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
		return errors.Wrap(err, "drop user table failed")
	}
	if err := db.Migrator().DropTable(&v1.Policy{}); err != nil {
		return errors.Wrap(err, "drop policy table failed")
	}
	if err := db.Migrator().DropTable(&v1.Secret{}); err != nil {
		return errors.Wrap(err, "drop secret table failed")
	}

	return nil
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
//
//nolint:unused // may be reused in the feature, or just show a migrate usage.
func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return errors.Wrap(err, "migrate user model failed")
	}
	if err := db.AutoMigrate(&v1.Policy{}); err != nil {
		return errors.Wrap(err, "migrate policy model failed")
	}
	if err := db.AutoMigrate(&v1.Secret{}); err != nil {
		return errors.Wrap(err, "migrate secret model failed")
	}

	return nil
}

// resetDatabase resets the database tables.
//
//nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
func resetDatabase(db *gorm.DB) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}

	return nil
}
