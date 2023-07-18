package ingame

import (
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/comminfo"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/http_client"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"sort"

	"github.com/sirupsen/logrus"
)

func RequestGameFTAssets(appid int, account string, uid uint64) (common.SortRGameServerERC20Assets, error) {
	type tmpErc20Ret struct {
		Code    int64                             `json:"code"`
		Message string                            `json:"message"`
		Data    common.SortRGameServerERC20Assets `json:"data"`
	}
	erc20Ret := tmpErc20Ret{
		Code:    0,
		Message: "",
		Data:    common.SortRGameServerERC20Assets{},
	}

	baseurl, err := comminfo.GetBaseUrl(appid)
	if err != nil {
		return nil, fmt.Errorf("RequestGameFTAssets GetBaseUrl failed, %v", err)
	}

	url := fmt.Sprintf("%s%s?appId=%d&uid=%d", baseurl, const_def.GameServer_FT_Query, appid, uid)

	err = http_client.HttpClientReqWithGet(url, &erc20Ret)
	logger.Logrus.WithFields(logrus.Fields{"error": err, "erc20ret": erc20Ret, "url": url}).Info("RequestGameFTAssets")
	if err != nil {

		return nil, fmt.Errorf("request failed ,%v", err)
	}

	if erc20Ret.Code != common.SuccessCode {

		return nil, fmt.Errorf(erc20Ret.Message)
	}

	sort.Sort(erc20Ret.Data)

	return erc20Ret.Data, nil
}

func RequestGameNFTAssets(appid int, account, assetName string, page, pageSize int64, uid uint64) ([]common.RGameServerERC721Asset, int64, error) {

	erc721Ret := common.RNFTAssetResponse{
		Code:    0,
		Message: "",
		Total:   0,
		Data:    []common.RGameServerERC721Asset{},
	}

	baseurl, err := comminfo.GetBaseUrl(appid)
	if err != nil {
		return nil, 0, fmt.Errorf("RequestGameNFTAssets GetBaseUrl failed, %v", err)
	}

	enurl := fmt.Sprintf("%s%s?appId=%d&uid=%d&assetName=%s&page=%d&pageSize=%d", baseurl, const_def.GameServer_NFT_Query, appid, uid, assetName, page, pageSize)

	err = http_client.HttpClientReqWithGet(enurl, &erc721Ret)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"error": err, "erc721ret": erc721Ret, "url": enurl}).Error("RequestGameNFTAssets request failed")
		return nil, 0, fmt.Errorf("RequestGameNFTAssets request failed ,%v", err)
	}

	logger.Logrus.WithFields(logrus.Fields{"erc721ret": erc721Ret, "url": enurl}).Info("RequestGameNFTAssets info")

	if erc721Ret.Code != common.SuccessCode || int(erc721Ret.Total) < len(erc721Ret.Data) || len(erc721Ret.Data) > int(pageSize) {
		return nil, 0, fmt.Errorf("RequestGameNFTAssets response status code failed ,%v", erc721Ret.Code)
	}

	return erc721Ret.Data, erc721Ret.Total, nil
}

func RequestGameNFTAssetDetail(appid int, equipID string, skipGameValidate bool) (*common.RGameServerERC721AssetDetail, error) {
	type tmpErc721Ret struct {
		Code    int64                               `json:"code"`
		Message string                              `json:"message"`
		Data    common.RGameServerERC721AssetDetail `json:"data"`
	}
	queryAssetDetailRet := tmpErc721Ret{
		Code:    0,
		Message: "",
		Data:    common.RGameServerERC721AssetDetail{},
	}

	baseurl, err := comminfo.GetBaseUrl(appid)
	if err != nil {
		return nil, fmt.Errorf("RequestGameNFTAssets GetBaseUrl failed, %v", err)
	}

	url := fmt.Sprintf("%s%s?appId=%d&equipment_id=%s", baseurl, const_def.GameServer_NFT_Detail_Query, appid, equipID)

	if skipGameValidate {
		url += "&skip_validate=1"
	}
	err = http_client.HttpClientReqWithGet(url, &queryAssetDetailRet)
	logger.Logrus.WithFields(logrus.Fields{"error": err, "detail": queryAssetDetailRet, "url": url}).Info("RequestGameNFTAssetDetail")
	if err != nil {

		return nil, fmt.Errorf("RequestGameNFTAssets request failed ,%v", err)
	}
	if queryAssetDetailRet.Code != common.SuccessCode {

		return nil, fmt.Errorf("response: %v", queryAssetDetailRet)
	}

	return &queryAssetDetailRet.Data, nil
}
