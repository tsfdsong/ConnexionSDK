package common

type IChainListener interface {
	Listen() error
	Stop()
	GetGameId() int
	HandleSpecHeight(int64) error
}
