package economia

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "hunt",
		Aliases:     []string{"cazar", "caza"},
		Description: "caza animales y gana cometas y XP",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Hunt,
	})
}

func Hunt(evt *events.Message, client *whatsmeow.Client, _ []string) {
	runNewEconomyActivity("hunt", evt, client)
}