package economia

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "dig",
		Aliases:     []string{"excavar", "excavacion"},
		Description: "excava en busca de materiales y cometas",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Dig,
	})
}

func Dig(evt *events.Message, client *whatsmeow.Client, _ []string) {
	runNewEconomyActivity("dig", evt, client)
}