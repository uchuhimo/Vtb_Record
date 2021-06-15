package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Config *MainConfig
var ConfigChanged bool

type CQJsonConfig struct {
	NeedCQBot bool
	QQGroupID []int
	CQHost    string
	CQToken   string
}
type UsersConfig struct {
	TargetId          string
	Name              string
	DownloadDir       string
	NeedDownload      bool
	StreamLinkArgs    []string
	AltStreamLinkArgs []string
	AltProxy          string
	CQConfig          CQJsonConfig
}
type ModuleConfig struct {
	//EnableProxy     bool
	//Proxy           string
	Name             string
	Enable           bool
	Users            []UsersConfig
	DownloadProvider string
	UseFollowPolling bool
	ApiHostUrl       string
	EnableProxy      bool
	Proxy            string
	Cookie           string
	PollInterval     int
}
type MainConfig struct {
	CriticalCheckSec int
	NormalCheckSec   int
	LogFile          string
	LogFileSize      int
	LogLevel         string
	RLogLevel        string
	DownloadQuality  string
	DownloadDir      []string
	UploadDir        string
	Module           []ModuleConfig
	DanmuHost        string
	DanmuPort        string
	EnableDanmu      bool
	EnableTS2MP4     bool
}

var v *viper.Viper

func InitConfig() {
	log.Print("Init config!")
	initConfig()
	log.Print("Load config!")
	_, _ = ReloadConfig()
	//fmt.Println(Config)
}

func initConfig() {
	/*viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.SetConfigType("json")*/
	v = viper.NewWithOptions(viper.KeyDelimiter("::::"))
	v.SetConfigFile(viper.ConfigFileUsed())
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		os.Exit(1)
	}

	ConfigChanged = true
	v.OnConfigChange(func(in fsnotify.Event) {
		ConfigChanged = true
	})
}

func ReloadConfig() (bool, error) {
	if !ConfigChanged {
		return false, nil
	}
	ConfigChanged = false
	err := v.ReadInConfig()
	if err != nil {
		return true, err
	}
	config := &MainConfig{}
	err = v.Unmarshal(config)
	if err != nil {
		fmt.Printf("Struct config error: %s", err)
	}
	Config = config
	UpdateLogLevel()
	return true, nil
}

func LevelStrParse(levelStr string) (level logrus.Level) {
	level = logrus.InfoLevel
	if levelStr == "debug" {
		level = logrus.DebugLevel
	} else if levelStr == "info" {
		level = logrus.InfoLevel
	} else if levelStr == "warn" {
		level = logrus.WarnLevel
	} else if levelStr == "error" {
		level = logrus.ErrorLevel
	}
	return level
}

func UpdateLogLevel() {
	if ConsoleHook != nil {
		level := LevelStrParse(Config.LogLevel)
		ConsoleHook.LogLevel = level
		logrus.Printf("Set logrus console level to %s", level)
	}
}

func PrepareConfig() {
	confPath := flag.String("config", "config.yaml", "config.yaml location")
	flag.Parse()
	viper.SetConfigFile(*confPath)
	InitConfig()
}
