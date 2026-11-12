package economia

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "einfo",
		Aliases:     []string{"economia"},
		Description: "muestra toda la información del sistema económico",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     EInfo,
	})
}

func formatEInfoDuration(seconds int64) string {
	if seconds < 0 {
		seconds = 0
	}

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	if secs > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}

	if len(parts) == 0 {
		return "0s"
	}

	return strings.Join(parts, " ")
}

func getEInfoCooldown(jid, action string, duration int64) string {
	if store.GlobalEconomy == nil {
		return "No disponible"
	}

	lastUse := store.GlobalEconomy.GetCooldown(jid, action)

	if lastUse <= 0 {
		return "Disponible"
	}

	remaining := duration - (time.Now().Unix() - lastUse)

	if remaining <= 0 {
		return "Disponible"
	}

	return formatEInfoDuration(remaining)
}

func appendEInfoCooldown(
	body *strings.Builder,
	jid string,
	command string,
	description string,
	action string,
	duration int64,
) {
	remaining := getEInfoCooldown(jid, action, duration)

	body.WriteString(fmt.Sprintf(
		"• `%s` — %s\n  ⏳ %s\n",
		command,
		description,
		remaining,
	))
}

func getShopItemNumber(itemID string) int {
	for index, item := range store.ShopItems {
		if item.ID == itemID {
			return index + 1
		}
	}

	return 0
}

func EInfo(evt *events.Message, client *whatsmeow.Client, _ []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	jid := evt.Info.Sender.String()

	balance := store.GlobalEconomy.GetBalance(jid)
	bank := store.GlobalEconomy.GetBank(jid)
	level, xp := store.GlobalEconomy.GetProfile(jid)
	neededXP := level * 100

	var body strings.Builder

	body.WriteString(
		"╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🪙 SISTEMA ECONÓMICO` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n",
	)

	body.WriteString("*📊 INFORMACIÓN PERSONAL*\n")
	body.WriteString(fmt.Sprintf("👛 *Billetera:* %d cometas\n", balance))
	body.WriteString(fmt.Sprintf("🏦 *Banco:* %d cometas\n", bank))
	body.WriteString(fmt.Sprintf("💰 *Patrimonio:* %d cometas\n", balance+bank))
	body.WriteString(fmt.Sprintf("⭐ *Nivel:* %d\n", level))
	body.WriteString(fmt.Sprintf("✨ *Experiencia:* %d/%d XP\n", xp, neededXP))

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*⚒️ ACTIVIDADES ECONÓMICAS*\n")
	body.WriteString("Consigue cometas y experiencia mediante diferentes actividades.\n\n")

	appendEInfoCooldown(&body, jid, "/w", "Trabajo", "work", 600)
	appendEInfoCooldown(&body, jid, "/adventure", "Aventura", "adventure", 1200)
	appendEInfoCooldown(&body, jid, "/mine", "Minería", "mine", 900)
	appendEInfoCooldown(&body, jid, "/fish", "Pesca", "fish", 900)
	appendEInfoCooldown(&body, jid, "/mendigar", "Mendigar", "beg", 300)
	appendEInfoCooldown(&body, jid, "/slut", "Actividad especial", "slut", 1200)
	appendEInfoCooldown(&body, jid, "/hunt", "Caza", "hunt", 600)
	appendEInfoCooldown(&body, jid, "/dig", "Excavación", "dig", 480)
	appendEInfoCooldown(&body, jid, "/talar", "Tala de árboles", "talar", 1500)
	appendEInfoCooldown(&body, jid, "/crime", "Crimen", "crime", 1800)

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🎁 RECOMPENSAS PERIÓDICAS*\n")

	appendEInfoCooldown(&body, jid, "/daily", "Recompensa diaria", "daily", 86400)
	appendEInfoCooldown(&body, jid, "/weekly", "Recompensa semanal", "weekly", 7*86400)
	appendEInfoCooldown(&body, jid, "/monthly", "Recompensa mensual", "monthly", 30*86400)

	body.WriteString("\n🔥 `/streak` — Consulta tu racha diaria.\n")

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🏦 BANCO Y ADMINISTRACIÓN*\n")
	body.WriteString("💰 `/bal` — Consulta tu billetera, banco y experiencia.\n")
	body.WriteString("📥 `/dep <cantidad>` — Deposita cometas.\n")
	body.WriteString("📤 `/with <cantidad>` — Retira cometas.\n")
	body.WriteString("💸 `/pay @usuario <cantidad>` — Transfiere cometas.\n")

	appendEInfoCooldown(&body, jid, "/rob @usuario", "Intentar robar cometas", "rob", 3600)

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🎯 MISIONES Y PROGRESIÓN*\n")
	body.WriteString("📋 `/mission` — Consulta tus misiones económicas.\n")
	body.WriteString("⭐ `/rank` — Consulta tu nivel y progreso.\n")
	body.WriteString("🏆 `/leveltop` — Ranking de niveles.\n")
	body.WriteString("💎 `/baltop` — Ranking de riqueza.\n")

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🛒 TIENDA, HERRAMIENTAS Y OBJETOS*\n")
	body.WriteString("🏪 `/shop` — Consulta todos los objetos disponibles.\n")
	body.WriteString("🛍️ `/buy <número>` — Compra un objeto o herramienta.\n")
	body.WriteString("🎒 `/ecoinv` — Consulta tu inventario económico.\n")
	body.WriteString("🧪 `/use <objeto>` — Activa objetos compatibles.\n")
	body.WriteString("💵 `/sellitem <objeto>` — Vende objetos de tu inventario.\n")
	body.WriteString("📦 `/chest` — Abre tus cofres.\n\n")

	body.WriteString("*Herramientas necesarias para algunas actividades:*\n")

	for _, equipment := range economyEquipment {
		item, ok := store.GetItem(equipment.ItemID)
		if !ok {
			continue
		}

		number := getShopItemNumber(equipment.ItemID)

		if number > 0 {
			body.WriteString(fmt.Sprintf(
				"• %s *%s* — `/%s` · %d cometas · `/buy %d`\n",
				equipment.Emoji,
				item.Name,
				equipment.Action,
				item.Price,
				number,
			))
		} else {
			body.WriteString(fmt.Sprintf(
				"• %s *%s* — necesario para `/%s`.\n",
				equipment.Emoji,
				item.Name,
				equipment.Action,
			))
		}
	}

	body.WriteString(
		"\nLas herramientas permanecen en tu inventario. No se consumen al realizar actividades.\n",
	)

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🎰 APUESTAS CLÁSICAS*\n")
	body.WriteString("🪙 `/flip` — Cara o cruz · ⏳ Sin cooldown.\n")
	body.WriteString("🎡 `/roulette` — Ruleta de colores y números · ⏳ Sin cooldown.\n")
	body.WriteString("🎰 `/slots` — Máquina tragamonedas · ⏳ Sin cooldown.\n")

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")

	body.WriteString("*🎲 JUEGOS DE APUESTAS*\n")
	body.WriteString("🎲 `/dice` — Apuesta con dados · ⏳ Sin cooldown.\n")
	body.WriteString("📈 `/highlow` — Alto o bajo · ⏳ Sin cooldown.\n")
	body.WriteString("✊ `/rpsbet` — Piedra, papel o tijera · ⏳ Sin cooldown.\n")
	body.WriteString("🃏 `/blackjack` — Blackjack · ⏳ Sin cooldown.\n")
	body.WriteString("🏇 `/race` — Carrera · ⏳ Sin cooldown.\n")
	body.WriteString("🔵 `/plinko` — Plinko · ⏳ Sin cooldown.\n")
	body.WriteString("📊 `/limbo` — Multiplicadores · ⏳ Sin cooldown.\n")
	body.WriteString("🔢 `/numberbet` — Apuesta numérica · ⏳ Sin cooldown.\n")
	body.WriteString("💥 `/crash` — Crash · ⏳ Sin cooldown.\n")
	body.WriteString("⚖️ `/evenodd` — Par o impar · ⏳ Sin cooldown.\n")
	body.WriteString("💣 `/minesbet` — Minas · ⏳ Sin cooldown.\n")
	body.WriteString("🎁 `/jackpot` — Bote mayor · ⏳ Sin cooldown.\n")
	body.WriteString("🎡 `/wheel` — Rueda de premios · ⏳ Sin cooldown.\n")
	body.WriteString("🃏 `/baccarat` — Baccarat · ⏳ Sin cooldown.\n")
	body.WriteString("⚔️ `/war` — Guerra de cartas · ⏳ Sin cooldown.\n")

	body.WriteString("\n📖 `/bethelp` — Ayuda e instrucciones de apuestas.\n")

	allwCooldown := getEInfoCooldown(jid, "allw", 3600)

	body.WriteString("\n━━━━━━━━━━━━━━━━━━\n")
	body.WriteString(fmt.Sprintf(
		"> Nota: también puedes usar `/allw` para ejecutar las diez actividades en conjunto. Cooldown global: %s. Requiere las herramientas indicadas arriba.\n",
		allwCooldown,
	))

	body.WriteString(
		"╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞",
	)

	client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(body.String()),
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: evt.Message,
			},
		},
	})
}