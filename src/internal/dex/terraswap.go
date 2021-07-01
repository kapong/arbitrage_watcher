package dex

import (
	"arbitrage_monitoring/internal/model"
	"arbitrage_monitoring/internal/networking"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	TerraswapEndpoint = "https://fcd.terra.dev"
)

type TerraswapResponse struct {
	Output TerraswapOutput `json:"result"`
}

type TerraswapOutput struct {
	ReturnAmount string `json:"return_amount"`
	SpreadAmount string `json:"spread_amount"`
	CommissionAmount string `json:"commission_amount"`
}

type Terraswap struct {
	DEX
}

func (t* Terraswap) GetPrice(swap *model.Swap) *model.Swap {
	
	var query string
	if strings.ToLower(swap.InputToken) == "bluna" { // except bluna is difference from others
		query = fmt.Sprintf(`%s/wasm/contracts/terra1jxazgm67et0ce260kvrpfv50acuushpjsz2y0p/store?query_msg={"simulation":{"offer_asset":{"amount":"%d","info":{"token":{"contract_addr":"%s"}}}}}`, 
			TerraswapEndpoint, 
			swap.InputAmount,
			swap.Contract,
		)
	} else {
		query = fmt.Sprintf(`%s/wasm/contracts/terra1jxazgm67et0ce260kvrpfv50acuushpjsz2y0p/store?query_msg={"simulation":{"offer_asset":{"amount":"%d","info":{"native_token":{"denom":"%s"}}}}}`, 
			TerraswapEndpoint, 
			swap.InputAmount, 
			strings.ToLower(swap.InputToken),
		)
	}
	
	data, err := networking.GetContent(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	response := &TerraswapResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		log.Println(err)
	}
	return t.parse(swap, &response.Output)
}

func (t* Terraswap) parse(swap *model.Swap, response *TerraswapOutput) *model.Swap {
	returnAmount, err := strconv.ParseUint(response.ReturnAmount, 0, 64)
	if err != nil {
		log.Println(err)
	}
	spreadAmount, err := strconv.ParseUint(response.SpreadAmount, 0, 64)
	if err != nil {
		log.Println(err)
	}
	commissionAmount, err := strconv.ParseUint(response.CommissionAmount, 0, 64)
	if err != nil {
		log.Println(err)
	}
	swap.ReturnAmount = returnAmount
	swap.SpreadAmount = spreadAmount
	swap.CommissionAmount = commissionAmount
	swap.Changed = float64(float64(returnAmount * 100) / float64(swap.InputAmount))-100
	return swap
}