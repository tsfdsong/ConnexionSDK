package state

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/commdata"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
)

//CheckAttrsExist check attributes is match
func CheckAttrsExist(conAddr, tokenID string, outer []commdata.EquipmentAttr) error {
	newAttrs, err := contracts.GetAttrsFromChain(conAddr, tokenID)
	if err != nil {
		return err
	}

	newValues := make(map[uint64]string, 0)
	for _, v := range newAttrs {
		key := v.AttributeID
		value := v.AttributeValue

		newValues[key] = value
	}

	isSame := true
	for _, v := range outer {
		oldlen := len(newValues)

		newValues[v.AttributeID] = v.AttributeValue

		newlen := len(newValues)
		if oldlen < newlen {
			isSame = false
			break
		}
	}

	if !isSame {
		return fmt.Errorf("outer attrs is not match contract attrs")
	}

	return nil
}
