package mysql

import (
	"example.com/m/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB  *gorm.DB
	err error
)

func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableAutomaticPing:                     true,
		PrepareStmt:                              false,
	})

	if err != nil {
		panic(err.Error())
		return err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("连接池失败:" + err.Error())
		panic(err.Error())
		return err
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(500)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(2000)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	////////////////////////////////////////////////////////////////
	model.CheckIsExistModelAdmin(DB)
	model.CheckIsExistModelConfig(DB)
	model.CheckIsExistModelMenu(DB)
	model.CheckIsExistModelRole(DB)
	model.CheckIsExistModelRoleMenu(DB)
	model.CheckIsExistModelPrepaidPhoneOrders(DB)
	model.CheckIsExistModelLog(DB)
	model.CheckIsExistModelReceiveAddress(DB)
	model.CheckIsExistModelBackLog(DB)
	model.CheckIsExistModelAccountChange(DB)
	model.CheckIsExistModelBalanceChange(DB)
	model.CheckIsExistModelConsoleManagementData(DB)
	////////////////////////////////////////////////////////////////

	return nil
}

func Close() {
	sqlDB, _ := DB.DB()
	sqlDB.Close()
}
