package main

import (
	"arbitrage_monitoring/internal/dex"
	"arbitrage_monitoring/internal/model"
	"arbitrage_monitoring/internal/networking"
	"arbitrage_monitoring/internal/utils"
	"arbitrage_monitoring/pkg/discord"
	"arbitrage_monitoring/pkg/webhook"
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

	webhook := webhook.Webhook{&discord.Discord{}}
	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK")
	discordLoggerWebhookURL := os.Getenv("DISCORD_LOGGER_WEBHOOK")

	fmt.Print("\n_____________________\n\n")
	fmt.Println("WATCHING_INTERVAL:", watchingInterval)
	fmt.Print("_____________________\n\n")

	terraswap := dex.DEX{&dex.Terraswap{}}

	latestSwap := make(map[int]float64)
	swapList := []*model.Swap{
		&model.Swap{ // Luna to Bluna
			Contract: "terra1jxazgm67et0ce260kvrpfv50acuushpjsz2y0p",
			InputToken: model.Token{
				Code: "uluna",
				Name: "🌕 LUNA",
			},
			OutputToken: model.Token{
				Code: "bluna",
				Name: "🌒 bLUNA",
			},
			InputAmount:         1000 * uint64(1000000),
			URL:                 "https://app.terraswap.io/#Swap",
			ExpectedProfit:      5.00,
			ExpectedProfileStep: 1,
		},
		&model.Swap{ // Bluna to Luna
			Contract: "terra1kc87mu460fwkqte29rquh4hc20m54fxwtsx7gp",
			InputToken: model.Token{
				Code: "bluna",
				Name: "🌒 bLUNA",
			},
			OutputToken: model.Token{
				Code: "uluna",
				Name: "🌕 LUNA",
			},
			InputAmount:         1000 * uint64(1000000),
			URL:                 "https://app.terraswap.io/#Swap",
			ExpectedProfit:      0,
			ExpectedProfileStep: 1,
		},
	}

	// notify configuration changes
	fmt.Println(webhook.GetConfigContent(swapList))
	networking.PostContent(
		discordWebhookURL,
		webhook.GetConfigContent(swapList),
	)

	for {
		swapResults := make([]*model.Swap, 0)
		for i, swap := range swapList {

			// get simulated price of each
			swapResult, err := terraswap.GetPrice(swap)
			if err != nil {
				log.Println(err)
				break
			}
			swapResults = append(swapResults, swapResult)

			// if simulated profit is greater than expected profit then notify to webhook
			if swapResult.Changed >= swapResult.ExpectedProfit {
				if (math.Abs(math.Abs(latestSwap[i])-math.Abs(swapResult.Changed)) > swapResult.ExpectedProfileStep) || latestSwap[i] == 0 {
					networking.PostContent(
						discordWebhookURL,
						webhook.GetContent(swapResult),
					)
					latestSwap[i] = swapResult.Changed
				}
			} else {
				latestSwap[i] = 0
			}
		}

		// arbitrager-logging channel
		networking.PostContent(
			discordLoggerWebhookURL,
			webhook.GetLoggingContent(swapResults),
		)

		time.Sleep(watchingInterval)
	}
}
