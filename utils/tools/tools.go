package tools

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

//use for random code
func GenCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	x := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(x)])
	}
	return sb.String()
}

func AsDashboardDisplayAmount(s string, displayNum int) (string, error) {

	ss := strings.Split(s, ".")
	if len(ss) != 1 && len(ss) != 2 {
		return "", errors.New("error amount")
	}

	if displayNum == 0 {
		return "", errors.New("error displayNum")
	}

	if s == "0" {
		return "0", nil
	}

	if len(ss) == 1 {
		return ss[0], nil
	} else {

		if ss[1] == strings.Repeat("0", len(ss[1])) {
			return ss[0], nil
		}

		displayLen := len(ss[1])
		if len(ss[1]) >= displayNum {
			displayLen = displayNum
		}

		fs := ss[1][:displayLen]

		if fs == strings.Repeat("0", displayLen) {
			return ss[0], nil
		} else {
			index := 0
			for i := displayLen - 1; i >= 0; i-- {
				if fs[i] == '0' {
					continue
				} else {
					index = i + 1
					break
				}
			}
			return (ss[0] + "." + fs[:index]), nil
		}
	}

}

func GetTokenAmount(s *big.Int, precision int32, displayNum int) (string, error) {
	if precision == 0 || displayNum == 0 {
		return "", errors.New("invalid args")
	}
	if s.Cmp(big.NewInt(0)) == 0 {
		return "0", nil
	}
	return GetTokenAmount2(s.String(), precision, displayNum)
}

func GetTokenAmount2(s string, precision int32, displayNum int) (string, error) {
	if precision == 0 || displayNum == 0 {
		return "", errors.New("invalid args")
	}

	if s == "0" {
		return "0", nil
	}

	if strings.HasPrefix(s, "0") {
		return "", errors.New("invalid args")
	}

	ft, _ := new(big.Int).SetString(s, 0)
	pr, _ := big.NewInt(0).SetString(fmt.Sprintf("1%s", strings.Repeat("0", int(precision))), 0)
	mod := big.NewInt(0)
	d := big.NewInt(0)
	d.Div(ft, pr)
	mod.Mod(ft, pr)

	m := strings.Repeat("0", int(precision)-len(mod.String())) + mod.String()
	exactAmount := fmt.Sprintf("%s.%s", d.String(), m)
	return AsDashboardDisplayAmount(exactAmount, displayNum)
}

func GetTokenExactAmount(s string, precision int32) (string, error) {
	if precision == 0 {
		return "", errors.New("invalid args")
	}

	if s == "0" {
		return "0", nil
	}

	if strings.HasPrefix(s, "0") {
		return "", errors.New("invalid args")
	}

	ft, _ := new(big.Int).SetString(s, 0)
	pr, _ := big.NewInt(0).SetString(fmt.Sprintf("1%s", strings.Repeat("0", int(precision))), 0)
	mod := big.NewInt(0)
	d := big.NewInt(0)
	d.Div(ft, pr)
	mod.Mod(ft, pr)
	//modLen := len(mod.String())

	if mod.Cmp(big.NewInt(0)) == 0 {
		return d.String(), nil
	} else {
		m := strings.Repeat("0", int(precision)-len(mod.String())) + mod.String()
		index := 0
		for i := len(m) - 1; i >= 0; i-- {
			if m[i] == '0' {
				continue
			} else {
				index = i + 1
				break
			}
		}
		return d.String() + "." + m[:index], nil
	}
}

//just for use at update db
//others require string > 0 && int > 0. if email/address(add tag.because address type has alias) check it
func CustomCheck(s interface{}) (map[string]interface{}, bool) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	r := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		tagKey := strings.ToLower(string(t.Field(i).Tag.Get("check")))
		value := v.Field(i)
		switch t.Field(i).Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
			if value.Int() != 0 {
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r[jsonKey] = value.Int()
			}

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			if value.Uint() != 0 {
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r[jsonKey] = value.Uint()
			}
		case reflect.Float32, reflect.Float64:
			if value.Float() != 0 {
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r[jsonKey] = value.Float()
			}

		case reflect.String:
			if value.String() != "" {
				if tagKey == "email" {
					if !CheckEmail(value.String()) {
						return r, false
					}
				} else if tagKey == "address" {
					if !CheckEthAddr(value.String()) {
						return r, false
					}
				}
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r[jsonKey] = value.String()
			}

		// case reflect.Bool:
		// 	r[key] = value.Bool()
		default:
		}
	}
	return r, true
}

func CheckEmail(s string) bool {
	emailRegexString := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	emailRegex := regexp.MustCompile(emailRegexString)
	return emailRegex.MatchString(s)
}

func CheckEthAddr(s string) bool {
	ethAddressRegexString := `^0x[0-9a-fA-F]{40}$`
	ethAddressRegex := regexp.MustCompile(ethAddressRegexString)
	return ethAddressRegex.MatchString(s)
}

//use for build query string

//key := string(t.Field(k).Tag.Get("json"))
func GetQueryConditionAndCheck(prefix string, s interface{}) (string, bool) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	r := ""
	for i := 0; i < t.NumField(); i++ {
		//key := t.Field(i).Name
		tagKey := strings.ToLower(string(t.Field(i).Tag.Get("check")))
		value := v.Field(i)
		switch t.Field(i).Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
			if value.Int() != 0 {
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r += fmt.Sprintf("%s%s=%v and ", prefix, jsonKey, value.Int())
			}

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			if value.Uint() != 0 {
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r += fmt.Sprintf("%s%s=%v and ", prefix, jsonKey, value.Uint())
			}

		case reflect.String:
			if value.String() != "" {
				if tagKey == "email" {
					if !CheckEmail(value.String()) {
						return r, false
					}
				} else if tagKey == "address" {
					if !CheckEthAddr(value.String()) {
						return r, false
					}
				}
				jsonKey := strings.ToLower(string(t.Field(i).Tag.Get("json")))
				r += fmt.Sprintf("%s%s=\"%v\" and ", prefix, jsonKey, value.String())
			}

		// case reflect.Bool:
		// 	r[key] = value.Bool()
		default:
		}
	}

	if len(r) > 0 {
		return r[:len(r)-4], true
	} else {
		return "", true
	}
}

//use for erc20deposit and nft mint with new..  only
//length is 2(erc20) or 3(nft mint with new)
//typeString order should be same as event arguments restrictly

func Unpack(hexstring string, typeString []string) ([]interface{}, error) {
	if len(typeString) != 2 && len(typeString) != 3 && len(typeString) != 4 && len(typeString) != 5 {
		return []interface{}{}, errors.New("invalid typeString")
	}

	if !strings.HasPrefix(hexstring, "0x") || len(hexstring) <= 3 {
		return []interface{}{}, errors.New("invalid hexstring")
	}

	data := hexstring[2:]

	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return []interface{}{}, errors.New("newtype failed")
	}
	addressTy, err := abi.NewType("address", "address", nil)
	if err != nil {
		return []interface{}{}, errors.New("newtype failed")
	}

	arguments := abi.Arguments{}
	for _, e := range typeString {
		if e == "uint256" {
			item := abi.Argument{
				Type: uint256Ty,
			}
			arguments = append(arguments, item)
		} else if e == "address" {
			item := abi.Argument{
				Type: addressTy,
			}
			arguments = append(arguments, item)
		} else {
			continue
		}
	}

	b, err := hex.DecodeString(data)
	if err != nil {
		return []interface{}{}, errors.New("decode hexstring failed")
	}
	s, err := arguments.Unpack(b)
	if err != nil {
		return []interface{}{}, errors.New("unpack failed")
	}
	return s, nil
}

func EthSignFix(hexs string) (string, error) {
	s := hexs
	f := hexs
	if strings.HasPrefix(hexs, "0x") {
		s = hexs[2:]
	} else {
		f = "0x" + hexs
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	if len(b) >= 65 {
		if b[64] == 0 || b[64] == 1 {
			b[64] += 27
		}
		return hexutil.Encode(b), nil
	} else {
		return f, nil
	}
}

func CheckGameDecimal(s string, gameDecimal, tokenDecimal int) (bool, string, error) {
	if s == "0" {
		return false, "", errors.New("error amount")
	}
	if tokenDecimal == 0 {
		return false, "", errors.New("error tokendecimal")
	}

	ss, err := GetTokenExactAmount(s, int32(tokenDecimal))
	if err != nil {
		return false, "", err
	}
	v := strings.Split(ss, ".")
	if len(v) != 1 && len(v) != 2 {
		return false, "", errors.New("error amount")
	}
	if len(v) == 1 {
		return true, v[0], nil
	} else {
		if len(v[1]) > gameDecimal {
			return false, "", errors.New("exceed tokendecimal")
		} else {
			return true, ss, nil
		}
	}
}

func GetTableID(s interface{}) string {
	defer func() {
		err := recover()
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err}).Error("s value have not ID ")
		}
	}()
	x := reflect.ValueOf(s).Elem().FieldByName("ID").Interface()
	return fmt.Sprintf("%+v", x)
}

func StringContains(src []string, dst string) bool {
	if len(src) == 0 {
		return false
	}

	for _, e := range src {
		if e == dst {
			return true
		}
	}
	return false
}

func GetBase64SVG(s string) (error, string) {
	firstBase64String := strings.Split(s, "data:application/json;base64,")
	if len(firstBase64String) != 2 {
		return fmt.Errorf("length:%v,firstBase64String:%+v", len(firstBase64String), firstBase64String), ""
	} else {
		firstDecodeString, err := base64.StdEncoding.DecodeString(firstBase64String[1])
		if err != nil {
			return err, ""
		}

		type tmpImage struct {
			Description string `json:"description"`
			Image       string `json:"image"`
		}

		t := tmpImage{}
		err = json.Unmarshal(firstDecodeString, &t)
		if err != nil {
			return err, ""
		}

		secondBase64String := strings.Split(t.Image, "data:image/svg+xml;base64,")
		if len(secondBase64String) != 2 {
			return fmt.Errorf("length:%v,secondBase64String:%+v", len(secondBase64String), secondBase64String), ""
		}

		secondDecodeString, err := base64.StdEncoding.DecodeString(secondBase64String[1])
		if err != nil {
			return err, ""
		}

		return nil, string(secondDecodeString)
	}
}

func CheckPersonalSign(jsonstring, signString, address string) bool {
	data := []byte(jsonstring)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	prefixData := crypto.Keccak256([]byte(msg))

	signature, err := hexutil.Decode(signString)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"err": err}).Info("decode sign failed")
		return false
	}
	if len(signature) != 65 {
		logger.Logrus.WithFields(logrus.Fields{"err": err}).Info("invalid sign")
		return false
	}
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(prefixData, signature)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"err": err}).Info("get publicKeyEDCSA failed")
		return false
	}

	addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()
	return (strings.ToLower(addr) == strings.ToLower(address))
}
