package main

import (
	"arbitrage_monitoring/internal/dex"
	"arbitrage_monitoring/internal/model"
	"arbitrage_monitoring/internal/networking"
	"arbitrage_monitoring/internal/utils"
	"arbitrage_monitoring/pkg/discord"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	utils.CheckEnv(
		"DISCORD_WEBHOOK",
		"WATCHING_INTERVAL",
	)

	watchingIntervalStr, err := strconv.ParseInt(os.Getenv("WATCHING_INTERVAL"), 10, 64)
	watchingInterval := time.Duration(watchingIntervalStr) * time.Minute
	utils.CheckErr(err)

	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK")

	fmt.Print("\n_____________________\n\n")
	fmt.Println("WATCHING_INTERVAL:", watchingInterval)
	fmt.Print("_____________________\n\n")

	terraswap := dex.DEX{&dex.Terraswap{}}

	latestSwap := make(map[int]float64, 0)
	swapList := []*model.Swap{
		&model.Swap{ // Luna to Bluna
			Contract: "terra1jxazgm67et0ce260kvrpfv50acuushpjsz2y0p",
			InputToken: model.Token{
				Code: "uluna",
				Name: "LUNA 🌕",
			},
			OutputToken: model.Token{
				Code: "bluna",
				Name: "bLUNA 🌒",
			},
			InputAmount:         1000 * uint64(1000000),
			URL:                 "https://app.terraswap.io/#Swap",
			ExpectedProfit:      5.00,
			ExpectedProfileStep: 5,
		},
		&model.Swap{ // Bluna to Luna
			Contract: "terra1kc87mu460fwkqte29rquh4hc20m54fxwtsx7gp",
			InputToken: model.Token{
				Code: "bluna",
				Name: "bLUNA 🌒",
			},
			OutputToken: model.Token{
				Code: "uluna",
				Name: "LUNA 🌕",
			},
			InputAmount:         1000 * uint64(1000000),
			URL:                 "https://app.terraswap.io/#Swap",
			ExpectedProfit:      -1.5,
			ExpectedProfileStep: 1.5,
		},
	}

	for {

		for i, swap := range swapList {
			// get simulated price of each
			swapResult := terraswap.GetPrice(swap)

			// if simulated profit is greater than expected profit then notify to webhook

			if swapResult.Changed >= swapResult.ExpectedProfit {
				if (math.Abs(math.Abs(latestSwap[i])-math.Abs(swapResult.Changed)) > swapResult.ExpectedProfileStep) || latestSwap[i] == 0 {
					networking.PostContent(
						discordWebhookURL,
						discord.GetDiscordContent(swapResult),
					)
					log.Println(fmt.Sprintf("%s → %s (%.2f%%) can be profitable !", swapResult.InputToken.Name, swapResult.OutputToken.Name, swapResult.Changed))
					latestSwap[i] = swapResult.Changed
				} else {
					log.Println(fmt.Sprintf("%s → %s (%.2f%%) is already emitted", swapResult.InputToken.Name, swapResult.OutputToken.Name, swapResult.Changed))
				}
			} else {
				log.Println(fmt.Sprintf("%s → %s (%.2f%%) doesn't met the threshold (expected >= %.2f%%)", swapResult.InputToken.Name, swapResult.OutputToken.Name, swapResult.Changed, swapResult.ExpectedProfit))
				latestSwap[i] = 0
			}

		}
		time.Sleep(watchingInterval)
	}
}
