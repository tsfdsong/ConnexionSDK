package main

import "github/Connector-Gamefi/ConnectorGoSDK/mock/game_server"

func main() {

	// configPath := flag.String("config_path", "/Users/hades/code/ConnectorGoSDK/cmd/", "config file")
	// logicLogFile := flag.String("logic_log_file", "/Users/hades/code/ConnectorGoSDK/cmd/log/sdk.log", "logic log file")
	// flag.Parse()

	// //init logic logger
	// logger.Init(*logicLogFile)
	// //load config
	// err := config.LoadConf(*configPath)
	// if err != nil {
	// 	logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("load config failed")
	// 	return
	// }
	// db := mysql.GetDB()
	// if db == nil {
	// 	fmt.Println("init db failed")
	// 	return
	// }

	// fmt.Println("hehe")
	// db_operation.Insert()
	game_server.GameServerMock()
}
