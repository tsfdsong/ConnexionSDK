package common

const (
	SuccessCode                = 0
	InnerError                 = 1001
	IncorrectParams            = 1002
	GameInfoNotExist           = 1003
	MustBindGameAccount        = 1004
	RequestGameServerFailed    = 1005
	UnsupportedNetwork         = 1006
	ContractAlreadyExist       = 1007
	ContractInfoNotExist       = 1008
	EquipmentNotExist          = 1009
	RecordNotExist             = 1010
	EmailCodeExist             = 1011
	WrongCode                  = 1012
	AlreayBound                = 1013
	NotBound                   = 1014
	FileSizeExceed             = 1015
	FailedUpload               = 1016
	UserNotExist               = 1017
	UploadToOSSFailed          = 1018
	GraphQueryFailed           = 1019
	OrderAlreadyExamine        = 1020
	MustBindEmail              = 1021
	WithdrawFailed             = 1022
	OrderNotSignYet            = 1023
	PrewithdrawCode            = 1024
	DepositNotiFailed          = 1025
	DepositNotiGameReject      = 1026
	WithdrawRecorverNotiFailed = 1027
	WithdrawCode               = 1028
	EmptySheet                 = 1029
	DuplicateAttrInfo          = 1030
	DepositSwitchClose         = 1031
	WithdrawSwitchClose        = 1032
	ParseLogSwitchClose        = 1033
	UnsupportedImageType       = 1034
	UserCantWithdraw           = 1035
	CantRepairOrder            = 1036
	CantClaimNow               = 1037
	GetTokenInfoFailed         = 1038
	GameUserNotExist           = 1039
	DepositMaxLimit            = 1040
	WithdrawBalanceNotEnough   = 1041
	FTRepeatClick              = 1042
	InvalidTimestamp           = 1043
	InvalidSignature           = 1044
	SendEmailFailed            = 1045
	NotFoundTx                 = 1046

	AuthCheck  = 2000
	AuthFailed = 2001
)

var ErrorMap = map[int]string{
	SuccessCode:                "success",
	InnerError:                 "inner error",
	IncorrectParams:            "Incorrect Parameters",
	GameInfoNotExist:           "Game Info Does Not Exist",
	MustBindGameAccount:        "Please Bind Game Account",
	RequestGameServerFailed:    "Request Game Server Failure",
	UnsupportedNetwork:         "Unsupported Network",
	ContractAlreadyExist:       "Contract Already Exist",
	ContractInfoNotExist:       "Missing Contract Info",
	EquipmentNotExist:          "Missing Gear Info",
	RecordNotExist:             "Record Does Not Exist",
	EmailCodeExist:             "Verification Code Has Been Sent",
	WrongCode:                  "Incorrect Verification Code",
	AlreayBound:                "Email Has Already Been Bound",
	NotBound:                   "Email Not Bound Yet",
	FileSizeExceed:             "File Size Exceeds Limit",
	FailedUpload:               "Upload Failure",
	UserNotExist:               "User Does Not Exist",
	UploadToOSSFailed:          "Upload to OSS Failure",
	GraphQueryFailed:           "Graph Query Failure",
	OrderAlreadyExamine:        "Order Does Not Need Examine",
	MustBindEmail:              "Please Bind Your Email",
	WithdrawFailed:             "Withdraw Failure",
	OrderNotSignYet:            "Withdraw Processing",
	PrewithdrawCode:            "Withdraw Failure",
	DepositNotiFailed:          "Replenish Order Failure",
	DepositNotiGameReject:      "Rejected Replenish Order",
	WithdrawRecorverNotiFailed: "Rejected Asset Recovery",
	WithdrawCode:               "Withdraw Failure",
	EmptySheet:                 "Empty Sheet",
	DuplicateAttrInfo:          "Duplicate Attribute Info",
	DepositSwitchClose:         "Paused Deposit",
	WithdrawSwitchClose:        "Paused Withdraw",
	ParseLogSwitchClose:        "Paused Parse Log",
	UnsupportedImageType:       "Unsupported Image Type",
	UserCantWithdraw:           "Your Withdraw is Paused",
	CantRepairOrder:            "Can not repair order",
	CantClaimNow:               "Withdraw Processing",
	AuthCheck:                  "Authorize validator failed",
	AuthFailed:                 "Authorithize Failed",
	GetTokenInfoFailed:         "Get Token Info Failed",
	GameUserNotExist:           "The account has not been registered in the game",
	DepositMaxLimit:            "deposit max limit",
	WithdrawBalanceNotEnough:   "game balance not enough",
	FTRepeatClick:              "Do not click repeatedly",
	InvalidTimestamp:           "invalid request",
	InvalidSignature:           "invalid signature",
	SendEmailFailed:            "send email failed",
	NotFoundTx:                 "not found tx",
}

type HpError struct {
	Inner error
	Code  int
	Msg   string
}

func (e *HpError) Error() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.Code != 0 {
		return ErrorMap[e.Code]
	} else {
		return "not defined error"
	}
}

func (e *HpError) CodeMsg() string {
	return e.Msg
}

func NewHpError(e error, code int, msg string) *HpError {
	if msg == "" {
		msg = ErrorMap[code]
	}

	return &HpError{
		Inner: e,
		Code:  code,
		Msg:   msg,
	}
}
