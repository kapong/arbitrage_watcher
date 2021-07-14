package discord

import (
	"arbitrage_monitoring/internal/model"
	"fmt"
)

func GetDiscordContent(swap *model.Swap) string {
	return fmt.Sprintf(`{"content":"> %s → %s\n*Expected profit*: **%.2f%%**\n%s"}`, swap.InputToken.Name, swap.OutputToken.Name, swap.Changed, swap.URL)
}
