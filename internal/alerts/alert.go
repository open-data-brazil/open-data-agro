package alerts

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Notifier sends operational alerts (webhook stub).
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a notifier. Empty webhookURL logs only.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: strings.TrimSpace(webhookURL),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify logs the message and POSTs to the webhook when configured.
func (n *Notifier) Notify(ctx context.Context, level slog.Level, message string, attrs ...any) {
	args := append([]any{"level", level.String(), "message", message}, attrs...)
	slog.Log(ctx, level, "ingestor alert", args...)

	if n.webhookURL == "" {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, strings.NewReader(message))
	if err != nil {
		slog.Error("alert webhook request failed", "error", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := n.client.Do(req)
	if err != nil {
		slog.Error("alert webhook delivery failed", "error", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
}
