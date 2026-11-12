package economia

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "crime",
		Aliases:     []string{"delito", "golpe"},
		Description: "intenta un golpe con riesgo de multa",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Crime,
	})
}

func Crime(evt *events.Message, client *whatsmeow.Client, _ []string) {
	runNewEconomyActivity("crime", evt, client)
}