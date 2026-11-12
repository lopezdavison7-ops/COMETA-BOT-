package economia

import (
	"context"
	"fmt"
	"strconv"
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
		Name:        "allw",
		Aliases:     []string{"workall", "trabajartodo", "todo"},
		Description: "ejecuta todas tus actividades económicas disponibles",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Allw,
	})
}

func Allw(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	results, err := store.GlobalEconomy.AllWork(evt.Info.Sender.String())
	if err != nil {
		if strings.HasPrefix(err.Error(), "cooldown:") {
			remaining, _ := strconv.ParseInt(
				strings.TrimPrefix(err.Error(), "cooldown:"),
				10,
				64,
			)

			sendEcoError(
				client,
				evt,
				fmt.Sprintf(
					"El paquete /allw está en cooldown. Volvé en *%d min %d seg*.",
					remaining/60,
					remaining%60,
				),
			)
			return
		}

		sendEcoError(client, evt, "No se pudieron ejecutar las actividades: "+err.Error())
		return
	}

	var body strings.Builder

	body.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `⚙️ ALL WORK` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

	total := 0
	totalXP := 0
	totalMissionBonus := 0
	executed := 0

	for _, result := range results {
		label := result.Title
		if label == "" {
			label = strings.ToUpper(result.Action)
		}

		if result.Skipped {
			body.WriteString(fmt.Sprintf(
				"⏳ *%s:* disponible en %d min %d seg\n",
				label,
				result.Remaining/60,
				result.Remaining%60,
			))
			continue
		}

		executed++
		total += result.Reward
		totalXP += result.XP
		totalMissionBonus += result.MissionBonus

		if result.Reward > 0 {
			body.WriteString(fmt.Sprintf(
				"✅ *%s:* %s (+%d)\n",
				label,
				result.Detail,
				result.Reward,
			))
		} else if result.Reward < 0 {
			body.WriteString(fmt.Sprintf(
				"⚠️ *%s:* %s (%d)\n",
				label,
				result.Detail,
				result.Reward,
			))
		} else {
			body.WriteString(fmt.Sprintf(
				"▫️ *%s:* %s\n",
				label,
				result.Detail,
			))
		}

		if result.XP > 0 && result.LeveledUp {
			body.WriteString(fmt.Sprintf(
				"   └ +%d XP · Nuevo nivel: %d\n",
				result.XP,
				result.Level,
			))
		}

		if result.MissionBonus > 0 {
			body.WriteString(fmt.Sprintf(
				"   └ Bonus de misión: +%d cometas\n",
				result.MissionBonus,
			))
		}
	}

	if executed == 0 {
		body.WriteString(
			"Todas tus actividades están en cooldown. No se activó el cooldown global de /allw.\n",
		)
	} else {
		body.WriteString(fmt.Sprintf(
			"\n💰 *Ganancia neta de actividades:* %d cometas\n",
			total,
		))

		if totalMissionBonus > 0 {
			body.WriteString(fmt.Sprintf(
				"🎯 *Premios de misiones:* +%d cometas\n",
				totalMissionBonus,
			))
		}

		if totalXP > 0 {
			body.WriteString(fmt.Sprintf(
				"✨ *Experiencia obtenida:* +%d XP\n",
				totalXP,
			))
		}

		body.WriteString(
			"⏳ Cooldown global: 60 minutos. Los cooldowns individuales se mantienen.\n",
		)
	}

	body.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(body.String()),
			ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
		},
	})
}