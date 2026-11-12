package economia

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "talar",
		Aliases:     []string{"woodcut", "cortar"},
		Description: "corta árboles y consigue cometas",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Talar,
	})
}

func Talar(evt *events.Message, client *whatsmeow.Client, _ []string) {
	runNewEconomyActivity("talar", evt, client)
}