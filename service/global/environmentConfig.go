package global

import (
	"fmt"
	"github.com/stevenroose/gonfig"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Params struct {
	Address              string `id:"address" short:"a" default:"0.0.0.0:2017" desc:"Listening address"`
	Config               string `id:"config" short:"c" default:"/etc/v2raya" desc:"v2rayA configuration directory"`
	V2rayBin             string `id:"v2ray-bin" desc:"Executable v2ray binary path. Auto-detect if put it empty."`
	V2rayConfigDirectory string `id:"v2ray-confdir" desc:"Additional v2ray config directory, files in it will be combined with config generated by v2rayA"`
	WebDir               string `id:"webdir" default:"/etc/v2raya/web" desc:"v2rayA web files directory"`
	Mode                 string `id:"mode" short:"m" desc:"(deprecated) Options: systemctl, service, universal. Auto-detect if not set"`
	PluginListenPort     int    `short:"s" default:"32346" desc:"ssr, pingTunnel, etc."`
	PassCheckRoot        bool   `desc:"Skip privilege checking. Use it only when you cannot start v2raya but confirm you have root privilege"`
	ResetPassword        bool   `id:"reset-password"`
	Verbose              bool   `id:"verbose" desc:"Mixedly print the log of v2rayA and v2ray-core"`
	ShowVersion          bool   `id:"version"`
}

var params Params

var dontLoadConfig bool

func initFunc() {
	defer SetServiceControlMode()
	if dontLoadConfig {
		return
	}
	err := gonfig.Load(&params, gonfig.Conf{
		FileDisable:       true,
		FlagIgnoreUnknown: false,
		EnvPrefix:         "V2RAYA_",
	})
	if err != nil {
		if err.Error() != "unexpected word while parsing flags: '-test.v'" {
			log.Fatal(err)
		}
	}
	// replace all dots of the filename with underlines
	params.Config = filepath.Join(
		filepath.Dir(params.Config),
		strings.ReplaceAll(filepath.Base(params.Config), ".", "_"),
	)
	if params.ShowVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
}

var once sync.Once

func GetEnvironmentConfig() *Params {
	once.Do(initFunc)
	return &params
}

func SetConfig(config Params) {
	params = config
}

func DontLoadConfig() {
	dontLoadConfig = true
}
