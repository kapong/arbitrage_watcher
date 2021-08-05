package webhook

import "arbitrage_monitoring/internal/model"

type WebhookInterface interface {
	GetConfigContent(swaps []*model.Swap) string
	GetContent(swap *model.Swap) string
	GetLoggingContent(swaps []*model.Swap) string
}

type Webhook struct {
	WebhookInterface
}
