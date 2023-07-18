package main

import (
	"flag"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	_ "github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/redis"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"github/Connector-Gamefi/ConnectorGoSDK/web"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	local := []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "run --Group=test --DataIds=config_name1,config_name2 --NacosAddrs=127.0.0.1:8848,127.0.0.1:8849 --NamespaceId=test --NacosLogLevel=debug",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "Group", Usage: "--Group"},
				&cli.StringFlag{Name: "DataIds", Usage: "--DataIds"},
				&cli.StringFlag{Name: "NacosAddrs", Usage: "--NacosAddrs"},
				&cli.StringFlag{Name: "NamespaceId", Usage: "--NamespaceId"},
				&cli.StringFlag{Name: "NacosLogLevel", Usage: "--NacosLogLevel"},
			},
			Action: func(cctx *cli.Context) error {
				run(cctx)
				return nil
			},
		},
	}
	app := &cli.App{
		Name:  "Connexion go_sdk server",
		Usage: "Connexion go_sdk server",

		Commands: local,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cctx *cli.Context) {
	env := os.Getenv("GO_ENV")
	configPath := flag.String("config_path", "./", "config file")
	logicLogFile := flag.String("logic_log_file", "./log/sdk.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	if env != "local" {
		nacosConf, err := config.NewNacosConfig(cctx)
		if err != nil {
			log.Fatal("Read config error:", err)
		}
		//load config
		err = config.LoadFromNacos(nacosConf)
		if err != nil {
			log.Fatal("Load nacos config error:", err)
		}
	} else {
		err := config.LoadConf(*configPath)
		if err != nil {
			log.Fatal("load config failed:", err)
		}
	}
	serverConf := config.GetServerConfig()
	if serverConf.LogOutStdout() {
		logger.Logrus.Out = os.Stdout
	}

	//set log level
	logger.SetLogLevel(serverConf.RunMode)

	db := mysql.GetDB()
	if db == nil {
		logger.Logrus.Error("init db failed")
		return
	}

	err := redis.InitRedis()
	if err != nil {
		logger.Logrus.Error("init redis failed")
		return
	}
	//init http pool
	pool.InitClient(int(config.GetHttpPoolSize()))
	pool.InitGraphClient(int(config.GetHttpPoolSize()))

	state.Loop()

	web.Run()
}
