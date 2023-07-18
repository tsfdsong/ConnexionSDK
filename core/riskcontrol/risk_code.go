package riskcontrol

const (
	CodeRiskNone     = 0 //init status
	CodeRiskWaiting  = 1 //waiting for review
	CodeRiskSuccess  = 2 //review success
	CodeRiskRejected = 3 // review rejected
	CodeRiskClosed   = 4 //risk control closed

	CodeRiskControlClose = 0 //close riskcontrol
	CodeRiskControlOpen  = 1 //open riskcontrol
)
