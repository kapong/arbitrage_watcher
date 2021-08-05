package discord

import (
	"arbitrage_monitoring/internal/model"
	"arbitrage_monitoring/pkg/webhook"
	"fmt"
)

type Discord struct {
	webhook.Webhook
}

func (d *Discord) GetConfigContent(swaps []*model.Swap) string {
	content := `🚨 **Arbitrage parameters has been changed** 🚨\n`
	for _, swap := range swaps {
		content = fmt.Sprintf(`%s → %s *(Expected %.2f%%, step %.2f%%)*\n`, content, swap.OutputToken.Name, swap.ExpectedProfit, swap.ExpectedProfileStep)
	}
	return fmt.Sprintf(`{"content":"%s"}`, content)
}

func (d *Discord) GetContent(swap *model.Swap) string {
	return fmt.Sprintf(`{"content":"> %s → %s\n*Expected profit*: **%.2f%%**\n%s"}`, swap.InputToken.Name, swap.OutputToken.Name, swap.Changed, swap.URL)
}

func (d *Discord) GetLoggingContent(swaps []*model.Swap) string {
	content := ""
	for _, swap := range swaps {
		content = fmt.Sprintf("%s → %s (%.2f%%)", content, swap.OutputToken.Name, swap.Changed)
	}
	return fmt.Sprintf(`{"content":"%s"}`, content)
}
