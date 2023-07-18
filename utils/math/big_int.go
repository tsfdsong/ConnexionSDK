package math

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"math/big"

	"github.com/shopspring/decimal"
)

// NewFromString convert string to big.Int
func NewFromString(in string) (*big.Int, error) {
	if in == "" {
		return big.NewInt(0), nil
	}

	data, ok := new(big.Int).SetString(in, 0)
	if !ok {
		return nil, fmt.Errorf("connot convert %s to big int", in)
	}

	return data, nil
}

func AddString(x *big.Int, data string) (*big.Int, error) {
	y, err := NewFromString(data)
	if err != nil {
		return nil, err
	}

	x = new(big.Int).Add(x, y)
	return x, nil
}

func CmpString(x, y string) (int, error) {
	xbig, err := NewFromString(x)
	if err != nil {
		return 0, err
	}

	ybig, err := NewFromString(y)
	if err != nil {
		return 0, err
	}

	res := xbig.Cmp(ybig)
	return res, nil
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CheckAmountValidator amount with decimal
func CheckAmountValidator(amount string, gameDecimal, tokenDecimal int) error {
	actAmount := ToWei(amount, 0) //amount with decimal

	if actAmount.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("amount can not is zero")
	}

	tlDecimal := ToWei("1", tokenDecimal)
	decl := ToWei("1", gameDecimal)
	minDeciml := new(big.Int).Div(tlDecimal, decl)
	if actAmount.Cmp(minDeciml) < 0 {
		return fmt.Errorf("amount is smaller")
	}

	maxAmount := ToWei(int64(const_def.FT_MAX_LIMIT), tokenDecimal) //max limit is not with decimal

	if actAmount.Cmp(maxAmount) > 0 {
		return fmt.Errorf("more than max limit %s %s", actAmount.String(), maxAmount.String())
	}

	return nil
}
