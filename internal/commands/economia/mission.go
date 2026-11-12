package economia

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "mission",
		Aliases:     []string{"misiones", "missions"},
		Description: "muestra tus misiones económicas diarias",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Mission,
	})
}

func Mission(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	missions, err := store.GlobalEconomy.GetDailyMissions(evt.Info.Sender.String())
	if err != nil {
		sendEcoError(client, evt, "No se pudieron consultar las misiones: "+err.Error())
		return
	}

	var body strings.Builder

	body.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🎯 MISIONES DIARIAS` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

	for _, mission := range missions {
		status := "🔄"
		if mission.Completed {
			status = "✅"
		}

		body.WriteString(fmt.Sprintf("%s *%s*\n", status, mission.Name))
		body.WriteString(fmt.Sprintf("   %s\n", mission.Description))
		body.WriteString(fmt.Sprintf(
			"   Progreso: *%d/%d* · Premio: *%d* cometas\n\n",
			mission.Progress,
			mission.Target,
			mission.Reward,
		))
	}

	body.WriteString(
		"Los premios se entregan automáticamente al completar cada misión. Las misiones se reinician cada día.\n",
	)
	body.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(body.String()),
			ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
		},
	})
}