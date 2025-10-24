package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewMySQLDB(cfg Config) (*gorm.DB, error) {
	DSNTemplate := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	
	masterDSN := fmt.Sprintf(DSNTemplate,
		cfg.Master.Username,
		cfg.Master.Password,
		cfg.Master.Host,
		cfg.Master.Port,
		cfg.Master.DBName,
	)

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		Logger:          logger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
		CreateBatchSize: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}
	var slavesDSNs []gorm.Dialector
	for _, slave := range cfg.Slaves {
		slaveDSN := fmt.Sprintf(DSNTemplate,
			slave.Username,
			slave.Password,
			slave.Host,
			slave.Port,
			slave.DBName,
		)
		slavesDSNs = append(slavesDSNs, mysql.Open(slaveDSN))
	}
	dbresolverCfg := dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(masterDSN)},
		Replicas: slavesDSNs,
		Policy: dbresolver.RandomPolicy{},
	}
	if err := db.Use(dbresolver.Register(dbresolverCfg)); err != nil {
		return nil, err
	}
	return db, nil
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}