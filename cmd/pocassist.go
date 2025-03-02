package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jweny/pocassist/api/routers"
	conf2 "github.com/jweny/pocassist/pkg/conf"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/logging"
	"github.com/jweny/pocassist/pkg/util"
	"github.com/jweny/pocassist/poc/rule"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"path/filepath"
)

func init() {
	PrintBanner()
}

func InitAll() {
	// config 必须最先加载
	conf2.Setup()
	logging.Setup()
	db.Setup()
	routers.Setup()
	util.Setup()
	rule.Setup()
}

// 使用viper 对配置热加载
func HotConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("cmd.HotConf, fail to get current path: %v", err)
	}
	// 配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(dir, "config.yaml")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	// watch 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		log.Println("config file changed:", e.Name)
		InitAll()
	})
}

func RunApp() {
	app := cli.NewApp()
	app.Name = "pocassist"
	app.Usage = "Online POC Framework Without Writing Code"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			// 后端端口
			Name:  "port",
			Aliases: []string{"p"},
			Value: "1231",
			Usage: "web server `PORT`",
		},
	}
	app.Action = RunServer

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("cli.RunApp err: %v", err)
		return
	}
}

func RunServer(c *cli.Context) error {
	InitAll()
	HotConf()
	port := c.String("port")
	routers.InitRouter(port)
	return nil
}

