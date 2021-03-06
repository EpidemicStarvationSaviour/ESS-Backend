// parse the conf/app.ini and store the variable in the single struct
package setting

import (
	"log"
	"path/filepath"

	"github.com/go-ini/ini"

	"time"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

type App struct {
	RuntimeRootPath string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
}

var AppSetting = &App{}

type Server struct {
	RunMode         string
	HttpPort        int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	JwtExpireTime   time.Duration
	CacheSize       int
	CacheExpireTime time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	DbName      string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Admin struct {
	Password string
	UserId   int
	Name     string
	Phone    string
}

var AdminSetting = &Admin{}

type Secret struct {
	JwtKey    string
	JwtIssuer string
	SaltA     string
	SaltB     string
}

var SecretSetting = &Secret{}

type GRPC struct {
	Enable  bool
	Host    string
	Timeout uint
}

var GRPCSetting = &GRPC{}

type Amap struct {
	Enable    bool
	WebAPIKey string
}

var AmapSetting = &Amap{}

// init the setting struct
func Setup(path string) {
	Cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("Fail to parse `%v` : %v", path, err)
	}

	//---------------- app config ----------------------
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'AppSetting': %v", err)
	}
	// change the '/' to '\' in windows env, and do nothing in Unix
	AppSetting.RuntimeRootPath = filepath.FromSlash(AppSetting.RuntimeRootPath)
	AppSetting.LogSavePath = filepath.FromSlash(AppSetting.LogSavePath)

	//---------------- server config ----------------------
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'ServerSetting': %v", err)
	}
	ServerSetting.ReadTimeout *= time.Second
	ServerSetting.WriteTimeout *= time.Second
	ServerSetting.JwtExpireTime *= time.Minute
	ServerSetting.CacheExpireTime *= time.Minute

	//---------------- database config ----------------------
	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'DatabaseSetting': %v", err)
	}

	//---------------- admin config ----------------------
	err = Cfg.Section("admin").MapTo(AdminSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'AdminSetting': %v", err)
	}

	//---------------- secret config ----------------------
	err = Cfg.Section("secret").MapTo(SecretSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'SecretSetting': %v", err)
	}

	//---------------- gRPC config ----------------------
	err = Cfg.Section("grpc").MapTo(GRPCSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'gRPCSetting': %v", err)
	}

	//---------------- amap config ----------------------
	err = Cfg.Section("amap").MapTo(AmapSetting)
	if err != nil {
		log.Fatalf("Fail to parse 'AmapSetting': %v", err)
	}
}
