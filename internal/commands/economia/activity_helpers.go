package economia

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/store"
)

func sendActivityMessage(client *whatsmeow.Client, evt *events.Message, result store.ActivityResult) {
	if result.Skipped {
		remaining := result.Remaining
		sendEcoError(client, evt, fmt.Sprintf(
			"⏳ Ya realizaste esta actividad. Esperá *%d min %d seg*.",
			remaining/60,
			remaining%60,
		))
		return
	}

	var body strings.Builder

	body.WriteString(fmt.Sprintf(
		"╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `%s %s` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n",
		result.Emoji,
		result.Title,
	))

	body.WriteString("🫆֟፝𑑓᮫̤٘𑑓̸̷ " + result.Detail + "\n")

	if result.Reward > 0 {
		body.WriteString(fmt.Sprintf("💰 *Ganaste:* %d cometas\n", result.Reward))
	} else if result.Reward < 0 {
		body.WriteString(fmt.Sprintf("💸 *Multa:* %d cometas\n", -result.Reward))
	} else {
		body.WriteString("💀 *No obtuviste cometas esta vez.*\n")
	}

	if result.XP > 0 {
		body.WriteString(fmt.Sprintf("✨ *XP:* +%d\n", result.XP))
	}

	if result.LeveledUp {
		body.WriteString(fmt.Sprintf("🎉 *Subiste al nivel %d.*\n", result.Level))
	}

	if result.MissionBonus > 0 {
		body.WriteString(fmt.Sprintf(
			"🎯 *Bonus de misiones:* +%d cometas\n",
			result.MissionBonus,
		))
	}

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

func runNewEconomyActivity(action string, evt *events.Message, client *whatsmeow.Client) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	result, err := store.GlobalEconomy.RunActivity(
		evt.Info.Sender.String(),
		action,
	)
	if err != nil {
		sendEcoError(client, evt, "No se pudo procesar la actividad: "+err.Error())
		return
	}

	sendActivityMessage(client, evt, result)
}