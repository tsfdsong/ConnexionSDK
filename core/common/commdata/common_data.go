package commdata

import (
	"encoding/json"
)

//equipment attribute
type EquipmentAttr struct {
	AttributeID    uint64 `json:"attribute_id"`
	AttributeValue string `json:"attribute_value"`
	AttributeName  string `json:"attribute_name"`
}

//NftSignatureSrcData source data of nft withdraw tx
type NftSignatureSrcData struct {
	AttrIDs          []uint64 `json:"attr_ids"`
	AttrValues       []string `json:"attr_values"`
	AttrIndexsUpdate []int    `json:"attr_indexs_update"`
	AttrValuesUpdate []string `json:"attr_values_update"`
	AttrIndexsRM     []int    `json:"attr_indexs_delete"`
}

//Serialize serilize src date to bytes
func (data *NftSignatureSrcData) MarshalJson() ([]byte, error) {
	return json.Marshal(&data)
}
