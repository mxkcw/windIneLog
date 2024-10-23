package WindIne_orm_mysql

import (
	"errors"
	"fmt"
	"github.com/mxkcw/windIneLog/windIne_log"
	"github.com/mxkcw/windIneLog/windIne_orm/windIne_orm_config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

type GTORMMysql struct {
	MysqlDB    *gorm.DB
	MysqlError error
	mux        sync.RWMutex
}

var (
	mysqlOnce     sync.Once
	mysqlInstance *GTORMMysql
)

func Instance() *GTORMMysql {
	mysqlOnce.Do(func() {
		mysqlInstance = &GTORMMysql{}
	})
	mysqlInstance.MysqlError = nil
	return mysqlInstance
}

func CloseDB() {
	db, err := mysqlInstance.MysqlDB.DB()
	if err != nil {
		windIne_log.LogErrorf("db err: %v", err)
	}
	if err = db.Close(); err != nil {
		windIne_log.LogErrorf("close db failed: %v", err)
	}
}

func (aMysql *GTORMMysql) OPenMysql(dbUser string, dbPwd string, dbName string, dbAddress string, dbPort int, timeZone windIne_orm_config.WindIneTimeZone, endFunc func(err error)) {
	aMysql.mux.Lock()
	defer aMysql.mux.Unlock()
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s", dbUser, dbPwd, dbAddress, dbPort, dbName, timeZone.String())
	alogleve := logger.Silent

	aMysql.MysqlDB, aMysql.MysqlError = gorm.Open(mysql.Open(connectionStr), &gorm.Config{
		Logger: logger.Default.LogMode(alogleve),
	})
	if aMysql.MysqlError != nil {
		endFunc(errors.New(fmt.Sprintf("连接数据库失败==%s", aMysql.MysqlError)))
	} else {
		sqlDb, _ := aMysql.MysqlDB.DB()
		sqlDb.SetMaxOpenConns(5)
		sqlDb.SetMaxIdleConns(2)
		sqlDb.SetConnMaxIdleTime(time.Minute)
		//fmt.Printf("数据库==%s,连接成功", dbName)
		windIne_log.LogInfof("数据库==%s,连接成功", dbName)
		endFunc(nil)
	}
}

// InsertData 单例
func (aMysql *GTORMMysql) InsertData(dataModel interface{}) error {
	aMysql.mux.Lock()
	defer aMysql.mux.Unlock()

	result := aMysql.MysqlDB.Where(dataModel).Limit(1).Find(dataModel)
	err := result.Error
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		if curer := aMysql.MysqlDB.Create(dataModel).Error; curer != nil {
			return curer
		}
	}
	return nil
}
