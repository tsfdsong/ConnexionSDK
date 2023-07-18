package comminfo

import (
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/contracts"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/tools"
	"strings"
)

func GetChainSVG(url, contractAddr, tokenId string) (string, error) {
	finalChainURI := ""
	if url != "" {
		if strings.Contains(url, "ipfs") {
			return config.GetBlindBoxImage(), nil
		}

		finalChainURI = url
	} else {
		err, rawURL := contracts.GetEquipmentTokenURI(contractAddr, tokenId)
		if err != nil {
			return "", err
		}

		finalChainURI = rawURL
	}

	if finalChainURI == "" {
		return "", nil
	}

	err, svg := tools.GetBase64SVG(finalChainURI)
	if err != nil {
		return "", err
	}
	return svg, nil
}
