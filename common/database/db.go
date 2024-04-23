package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// sqlDatabase struct has necessary component for GORM configuration
type sqlDatabase struct {
	db     *gorm.DB
	config *Config
}

// SQLDatabase defines an interface to interact with DB
type SQLDatabase interface {
	Instance() *gorm.DB
	Connect() error
	Disconnect() error
}

// NewSQLDatabase will init a new instance of GORM based on the
// given configuration
func NewSQLDatabase(config *Config) (SQLDatabase, error) {
	db := &sqlDatabase{
		config: config,
	}
	err := db.Connect()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Instance will return the DB instance
func (db *sqlDatabase) Instance() *gorm.DB {
	return db.db
}

// Connect will connection with the database driver via GORM
func (db *sqlDatabase) Connect() error {
	var dialector gorm.Dialector
	var err error

	switch db.config.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.config.Username, db.config.Password, db.config.ConnectionString, db.config.Port, db.config.DBName)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			db.config.ConnectionString, db.config.Port, db.config.Username, db.config.DBName, db.config.Password)
		dialector = postgres.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			db.config.Username, db.config.Password, db.config.ConnectionString, db.config.Port, db.config.DBName)
		dialector = sqlserver.Open(dsn)
	default:
		return fmt.Errorf("unsupported database driver: %s", db.config.Driver)
	}

	db.db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(db.config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(db.config.MaxIdleConnections)

	return nil
}

// Disconnect will Disconnect from DB instance
func (db *sqlDatabase) Disconnect() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
