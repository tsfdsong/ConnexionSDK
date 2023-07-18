package state_test

import (
	"context"
	"flag"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/state"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/pool"
	"testing"
)

func TestFTSignHandler(t *testing.T) {
	configPath := flag.Arg(0)
	logicLogFile := flag.Arg(1)
	flag.Parse()
	t.Log("config_path:", configPath)
	t.Log("log_file:", logicLogFile)
	//init logic logger
	logger.Init(logicLogFile)

	err := config.LoadConf(configPath)
	if err != nil {
		t.Error(err)
		return
	}
	pool.InitClient(int(config.GetHttpPoolSize()))

	sig, err := state.FTSignHandler("500000000000000000", "0xb33851673a9a8fc83cd5879d7596d593cdc11f62", "0x1c4164255fbaad807d474cf154721feb99d3c037", "0x8c1a9cecda1414a536c2f9708bee843c9f59c4c2", "111288990978160433402525536557166674793271632763905433126106717502234129398454", 97, context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("sigCode:", sig.Code)
	t.Log("sigHash:", sig.ReqHash)
}
