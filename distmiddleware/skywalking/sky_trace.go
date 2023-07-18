package skywalking

import (
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"sync"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/sirupsen/logrus"
)

var tracer *go2sky.Tracer
var once sync.Once

func GetSkyTrace() *go2sky.Tracer {
	once.Do(func() {
		var err error
		rp, err := reporter.NewGRPCReporter(config.GetSkyWalkingConfig().SkyRPCNode, reporter.WithCheckInterval(5*time.Second))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Fatal("GetSkyTrace NewGRPCReporter failed")
			panic(err)
		}

		tracer, err = go2sky.NewTracer(config.GetSkyWalkingConfig().ServerName, go2sky.WithReporter(rp), go2sky.WithCorrelation(10, 256))
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Fatal("GetSkyTrace NewTracer failed")
			panic(err)
		}

		if tracer == nil {

			panic("instance is null")
		}
	})

	return tracer
}
