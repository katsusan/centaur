package config

import (
	"fmt"
	"io/ioutil"

	"github.com/Katsusan/centaur/internal/util"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type Config struct {
	db   *gorm.DB
	parm *Params
}

//Params stores the application config
type Params struct {
	Name           string
	Version        string
	Debug          bool   `yaml:"debug" flag:"debug"`
	LogLevel       string `yaml:"log-level" flag:"log-level"`
	ConfigFile     string
	ConfigPath     string `yaml:"config-path" flag:"config-path"`
	DbServerHost   string `yaml:"db-host" flag:"db-host"`
	DbServerPort   uint16 `yaml:"db-port" flag:"db-port"`
	DbUserName     string `yaml:"db-user" flag:"db-user"`
	DbPassword     string `yaml:"db-password" flag:"db-password"`
	DbName         string `yaml:"db-name" flag:"db-name"`
	HttpServerHost string `yaml:"http-host" flag:"http-host"`
	HttpServerPort uint16 `yaml:"http-port" flag:"http-port"`
	DaemonMode     bool   `yaml:"daemon" flag:"daemon"`
	DaemonPidPath  string `yaml:"daemon-pid-path" flag:"daemon-pid-path"`
	DaemonLogPath  string `yaml:"daemon-log-path" flag:"daemon-log-path"`
}

func initLogger(debug bool) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}

//NewParams will return an instance of struct Params
func NewParams(ctx *cli.Context) *Params {
	params := &Params{}
	params.Name = ctx.App.Name
	params.Version = ctx.App.Version

	return params
}

//NewConfig will return a Config instance with parameters generated from cli.Context
func NewConfig(ctx *cli.Context) *Config {
	initLogger(ctx.GlobalBool(("debug")))

	cfg := &Config{
		parm: NewParams((ctx)),
	}

	return cfg
}

//connectToDB will establish a new connection to mysql.
func (c *Config) connectToDB() error {
	DbDSN := c.DatabaseDSN()
	log.Debugln("will connect to ", DbDSN)

	db, err := gorm.Open("mysql", DbDSN)
	if err != nil {
		log.Fatal(err)
	}

	c.db = db
	return nil
}

//DatabaseDSN will return the connect string of database
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		c.parm.DbUserName,
		c.parm.DbPassword,
		c.parm.DbServerHost,
		c.parm.DbServerPort,
		c.parm.DbName,
	)
}

//LoadFromFile will read config from config file and parse it into *Params
func (p *Params) LoadFromFile(file string) error {
	if !util.FileExist(file) {
		return fmt.Errorf("config file[%s] not found", file)
	}

	yamlconfig, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlconfig, p)

}