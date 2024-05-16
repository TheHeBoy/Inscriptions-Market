package database

import (
	"database/sql"
	"fmt"
	"gohub/internal/model"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/plugin/dbresolver"

	"moul.io/zapgorm2"
	"time"

	"gorm.io/gorm"
)

// DB 对象
var DB *gorm.DB
var SQLDB *sql.DB

// SetupDB 初始化数据库和 ORM
func SetupDB() {

	var dbConfig gorm.Dialector
	switch config.Get("database.connection") {
	case "mysql":
		masterDb := config.GetStringMapString("database.mysql.master")
		dbConfig = mysql.New(mysql.Config{
			DSN: createDsn(masterDb),
		})
	default:
		panic("database connection not supported")
	}

	// 连接数据库
	err := Connect(dbConfig)
	if err != nil {
		panic(err)
	}

	maxOpenConns := config.GetInt("database.mysql.max_open_connections")
	maxIdleConns := config.GetInt("database.mysql.max_idle_connections")
	maxLifeSeconds := config.GetInt("database.mysql.max_life_seconds")

	SQLDB.SetMaxOpenConns(maxOpenConns)
	SQLDB.SetMaxIdleConns(maxIdleConns)
	SQLDB.SetConnMaxLifetime(time.Duration(maxLifeSeconds) * time.Second)

	// 添加 indexer 数据源
	indexerDb := config.GetStringMapString("database.mysql.indexer")
	err = DB.Use(
		dbresolver.Register(
			dbresolver.Config{Sources: []gorm.Dialector{mysql.Open(createDsn(indexerDb))}}, &model.HolderDO{}, &model.InscriptionDO{}, &model.Msc20DO{}, &model.ListDO{}, &model.TokenDO{}, "indexer").
			SetConnMaxLifetime(time.Duration(maxLifeSeconds) * time.Second).
			SetMaxIdleConns(maxIdleConns).
			SetMaxOpenConns(maxOpenConns),
	)
	if err != nil {
		panic(err)
	}

	err = DB.Migrator().AutoMigrate(&model.OrderDO{}, &model.OrderLogDO{})
	if err != nil {
		panic(err)
	}
}

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector) error {

	log := zapgorm2.New(logger.LogZap)
	log.IgnoreRecordNotFoundError = true
	log.SetAsDefault()

	// 使用 gorm.Open 连接数据库
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{Logger: log})

	// 处理错误
	if err != nil {
		return err
	}

	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		return err
	}
	return nil
}

func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() error {
	var err error
	switch config.Get("database.connection") {
	case "mysql":
		err = deleteMySQLTables()
	case "sqlite":
		err := deleteAllSqliteTables()
		if err != nil {
			return err
		}
	default:
		panic("database connection not supported")
	}

	return err
}

func deleteAllSqliteTables() error {
	tables := []string{}

	// 读取所有数据表
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteMySQLTables() error {
	dbname := CurrentDatabase()
	tables := []string{}

	// 读取所有数据表
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	DB.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	DB.Exec("SET foreign_key_checks = 1;")
	return nil
}

func TableName(obj interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	stmt.Parse(obj)
	return stmt.Schema.Table
}

func createDsn(data map[string]string) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		data["username"],
		data["password"],
		data["host"],
		data["port"],
		data["database"],
		data["charset"])
}
