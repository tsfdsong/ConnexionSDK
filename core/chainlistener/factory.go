package chainlistener

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/chainlistener/common"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
)

type Factory struct {
	Listeners []common.IChainListener
}

func NewListener(c config.FilterConfig, gameId int) (common.IChainListener, error) {
	return nil, errors.New("chain type not supported")
}

func NewFactory(c map[int]config.GameConfig) (*Factory, error) {
	f := &Factory{}
	for k, v := range c {
		if v.FilterConfig.ListenerType == "" {
			continue
		}
		listener, err := NewListener(v.FilterConfig, k)
		if err != nil {
			return nil, err
		}
		f.Listeners = append(f.Listeners, listener)
	}
	return f, nil
}

func (f *Factory) Listen() error {
	for _, listener := range f.Listeners {
		err := listener.Listen()
		if err != nil {
			return err
		}
		logger.Logrus.Info(fmt.Sprintf("start AppId:%d chain listener success", listener.GetGameId()))
	}
	return nil
}

func (f *Factory) Stop() {
	for _, listener := range f.Listeners {
		listener.Stop()
	}
}

func (f *Factory) Reload() error {
	f.Stop()
	f.Listeners = make([]common.IChainListener, 0)
	gameConfs := config.GetGameConfig()
	for k, v := range gameConfs {
		if v.FilterConfig.ListenerType == "" {
			continue
		}
		listener, err := NewListener(v.FilterConfig, k)
		if err != nil {
			return err
		}
		f.Listeners = append(f.Listeners, listener)
	}
	return nil
}

func (f *Factory) Restart() error {
	err := f.Reload()
	if err != nil {
		return err
	}
	return f.Listen()

}

func (f *Factory) OnConfigChange(callback func()) {
	channel := make(chan bool)
	config.RegistConfChange(channel)
	for range channel {
		callback()
	}
}
