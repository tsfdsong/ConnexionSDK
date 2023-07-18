package dashdata

type MultiAssetDeposit interface {
	Deposit() error
}

type MultiAssetWithdraw interface {
	Withdraw() error
}

type MultiAssetStrategy struct {
	depositImpl  MultiAssetDeposit
	withdrawImpl MultiAssetWithdraw
}

func (m *MultiAssetStrategy) SetDepositImpl(i MultiAssetDeposit) {
	m.depositImpl = i
}

func (m *MultiAssetStrategy) SetWithdrawImpl(i MultiAssetWithdraw) {
	m.withdrawImpl = i
}

func (m *MultiAssetStrategy) Deposit() error {
	return m.depositImpl.Deposit()
}

func (m *MultiAssetStrategy) Withdraw() error {
	return m.withdrawImpl.Withdraw()
}
