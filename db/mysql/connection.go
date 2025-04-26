package mysql

import (
	"log"
	"managedata/config"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	Db  *sqlx.DB
	err error
)

func ConnectMysqlDB() {
	cfg := mysql.Config{
		User:                 config.AppConfig.MysqlConf.Username,
		Passwd:               config.AppConfig.MysqlConf.Password,
		Net:                  config.AppConfig.MysqlConf.Net,
		Addr:                 config.AppConfig.MysqlConf.Address,
		DBName:               config.AppConfig.MysqlConf.DatabaseName,
		AllowNativePasswords: true,
	}

	// create a db connection pool
	Db, err = sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Error while connecting to MySQL DB: %v", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatalf("Error while connecting to MySQL DB: %v", err)
	}
	log.Println("Mysql Database Connected")
}

func CloseMysql() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Println("Error closing Mysql connection: ", err)
			return
		}
		log.Println("Closed Mysql connection Successfully!")
	}
}
