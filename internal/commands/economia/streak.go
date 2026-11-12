package economia

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "streak",
		Aliases:     []string{"racha", "daily_streak"},
		Description: "consulta tus días consecutivos reclamando /daily",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Streak,
	})
}

func Streak(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		return
	}

	streak, lastDate := store.GlobalEconomy.GetDailyStreak(evt.Info.Sender.String())

	text := fmt.Sprintf(
		"╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🔥 RACHA DIARIA` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n🔥 *Días consecutivos:* %d\n📅 *Último /daily:* %s\n\nReclamá /daily cada día para mantener tu racha.\n╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞",
		streak,
		lastDate,
	)

	if lastDate == "" {
		text = "Todavía no tenés una racha registrada. Usá /daily para comenzar."
	}

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
		},
	})
}