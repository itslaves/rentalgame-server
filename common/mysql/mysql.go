package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var mysqlClient *gorm.DB

func Init() error {
	if mysqlClient != nil {
		return errors.New("mysqlClient has already been initialized")
	}

	// https://godoc.org/github.com/go-sql-driver/mysql#Config
	params := make(map[string]string)
	for k, v := range viper.GetStringMap("mysql.params") {
		params[k] = fmt.Sprintf("%v", v)
	}
	var (
		err          error
		loc          *time.Location
		dialTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
	)
	loc, err = time.LoadLocation(viper.GetString("mysql.location"))
	if err != nil {
		return err
	}
	dialTimeout, err = time.ParseDuration(viper.GetString("mysql.dialTimeout"))
	if err != nil {
		return err
	}
	readTimeout, err = time.ParseDuration(viper.GetString("mysql.readTimeout"))
	if err != nil {
		return err
	}
	writeTimeout, err = time.ParseDuration(viper.GetString("mysql.writeTimeout"))
	if err != nil {
		return err
	}
	config := &mysql.Config{
		User:         viper.GetString("mysql.user"),
		Passwd:       viper.GetString("mysql_passwd"),
		Net:          "tcp",
		Addr:         viper.GetString("mysql.addr"),
		DBName:       viper.GetString("mysql.db"),
		Params:       params,
		Collation:    viper.GetString("mysql.collation"),
		Loc:          loc,
		Timeout:      dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ParseTime:    true,
	}

	mysqlClient, err = gorm.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return mysqlClient.Close()
}

func Client() *gorm.DB {
	return mysqlClient
}
