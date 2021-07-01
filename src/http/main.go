package main

import (
	"arbitrage_monitoring/internal/dex"
	"arbitrage_monitoring/internal/model"
	"arbitrage_monitoring/internal/networking"
	"arbitrage_monitoring/internal/utils"
	"arbitrage_monitoring/pkg/discord"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main(){
	utils.CheckEnv(
		"DISCORD_WEBHOOK",
		"EXPECTED_PROFIT",
		"WATCHING_INTERVAL",
	)

	expectedProfitThreshold, err := strconv.ParseFloat(os.Getenv("EXPECTED_PROFIT"), 64)
	utils.CheckErr(err)

	watchingIntervalStr, err := strconv.ParseInt(os.Getenv("WATCHING_INTERVAL"), 10, 64)
	watchingInterval := time.Duration(watchingIntervalStr) * time.Minute
	utils.CheckErr(err)

	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK")

	fmt.Print("\n_____________________\n\n")
	fmt.Println("EXPECTED_PROFIT:", expectedProfitThreshold)
	fmt.Println("WATCHING_INTERVAL:", watchingInterval)
	fmt.Print("_____________________\n\n")

	terraswap := dex.DEX{&dex.Terraswap{}}

	latestSwap := make(map[int]float64, 0)
	swapList := []*model.Swap{
		&model.Swap{ // Luna to Bluna
			Contract: "terra1jxazgm67et0ce260kvrpfv50acuushpjsz2y0p", 
			InputToken: "uLuna",
			OutputToken: "bLuna",
			InputAmount: 1000 * uint64(1000000),
			URL: "https://app.terraswap.io/#Swap",
		},
		&model.Swap{ // Bluna to Luna
			Contract: "terra1kc87mu460fwkqte29rquh4hc20m54fxwtsx7gp", 
			InputToken: "bLuna",
			OutputToken: "uLuna",
			InputAmount: 1000 * uint64(1000000),
			URL: "https://app.terraswap.io/#Swap",
		},
	}
	
	for {

		for i, swap := range swapList {
			// get simulated price of each
			swapResult := terraswap.GetPrice(swap)

			// if simulated profit is greater than expected profit then notify to webhook
			if swapResult.Changed >= expectedProfitThreshold {
				
				if swapResult.Changed > latestSwap[i] + expectedProfitThreshold {
					networking.PostContent(
						discordWebhookURL, 
						discord.GetDiscordContent(swapResult),
					)
					latestSwap[i] = swapResult.Changed
				} else {
					log.Println(fmt.Sprintf("%s → %s (%.2f%%) is already emitted", swapResult.InputToken, swapResult.OutputToken, swapResult.Changed))
				}
			} else {
				log.Println(fmt.Sprintf("%s → %s (%.2f%%) doesn't met the threshold", swapResult.InputToken, swapResult.OutputToken, swapResult.Changed))
				latestSwap[i] = 0
			}
			
		}
		time.Sleep(watchingInterval)
	}
}