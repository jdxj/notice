package main

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jdxj/notice/config"
)

var (
	db     *gorm.DB
	logger *zap.SugaredLogger
)

func init() {
	dbConfig := config.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}
