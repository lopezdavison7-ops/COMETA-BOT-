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

func getMentionedJID(evt *events.Message) string {
	if evt.Message == nil {
		return ""
	}

	if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
		if ctx := ext.GetContextInfo(); ctx != nil {
			mentions := ctx.GetMentionedJID()
			if len(mentions) > 0 {
				return mentions[0]
			}
		}
	}

	if img := evt.Message.GetImageMessage(); img != nil {
		if ctx := img.GetContextInfo(); ctx != nil {
			mentions := ctx.GetMentionedJID()
			if len(mentions) > 0 {
				return mentions[0]
			}
		}
	}

	return ""
}

func init() {
	registry.Register(&registry.Command{
		Name:        "rob",
		Aliases:     []string{"robar", "asaltar"},
		Description: "intenta robarle cometas a otro usuario",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Rob,
	})
}

func Rob(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		return
	}

	victimJID := getMentionedJID(evt)

	if victimJID == "" {
		client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String("❌ Debés mencionar a alguien para robarle. Ejemplo: `/rob @usuario`"),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: evt.Message,
				},
			},
		})
		return
	}

	thiefJID := evt.Info.Sender.String()

	if thiefJID == victimJID {
		client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String("❌ No te podés robar a vos mismo, chacho. 😂"),
				ContextInfo: &waE2E.ContextInfo{
					QuotedMessage: evt.Message,
				},
			},
		})
		return
	}

	var msg string
	var err error

	if remaining, ok := store.GlobalEconomy.CheckCooldown(thiefJID, "rob", 3600); !ok {
		msg = fmt.Sprintf(
			"La policía te está vigilando. Intentá robar de nuevo en %d minutos.",
			remaining/60,
		)
		err = fmt.Errorf("cooldown")
	} else if store.GlobalEconomy.ConsumeEffect(victimJID, "shield") {
		store.GlobalEconomy.SetCooldown(thiefJID, "rob")
		msg = "El escudo protector bloqueó el intento de robo. El escudo se consumió. 🛡️"
		err = fmt.Errorf("shield")
	} else {
		msg, err = store.GlobalEconomy.Rob(thiefJID, victimJID)

		if err == nil {
			store.GlobalEconomy.SetCooldown(thiefJID, "rob")
		}
	}

	var prefix string

	if err != nil {
		switch err.Error() {
		case "cooldown":
			prefix = "⏳ "
		case "poor":
			prefix = "😂 "
		case "shield":
			prefix = "🛡️ "
		default:
			prefix = "👮‍♂️ "
		}
	} else {
		prefix = "🥷 "
	}

	var body strings.Builder

	body.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🥷 ROB` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")
	body.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + prefix + msg + "\n")
	body.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(body.String()),
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: evt.Message,
			},
		},
	})
}