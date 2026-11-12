package economia

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "use",
		Aliases:     []string{"usar", "usaritem"},
		Description: "activa un objeto de tu inventario",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Use,
	})
}

func Use(evt *events.Message, client *whatsmeow.Client, args []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	if len(args) == 0 {
		sendEcoError(client, evt, "Usá: /use <número|id del objeto>.")
		return
	}

	jid := evt.Info.Sender.String()
	itemID := resolveInventoryItem(jid, strings.Join(args, " "))

	if itemID == "" {
		sendEcoError(client, evt, "No encontré ese objeto en tu inventario. Revisá /ecoinv.")
		return
	}

	switch itemID {
	case "xp_boost", "shield":
		if store.GlobalEconomy.HasEffect(jid, itemID) {
			if itemID == "xp_boost" {
				sendEcoError(client, evt, "Ya tenés una poción de XP activa.")
			} else {
				sendEcoError(client, evt, "Ya tenés un escudo activo.")
			}
			return
		}

		if !store.UseItem(jid, itemID) {
			sendEcoError(client, evt, "No se pudo retirar el objeto del inventario.")
			return
		}

		if err := store.GlobalEconomy.ActivateEffect(jid, itemID); err != nil {
			sendEcoError(
				client,
				evt,
				"El objeto se retiró, pero no se pudo activar: "+err.Error(),
			)
			return
		}

		if itemID == "xp_boost" {
			sendEcoSuccess(
				client,
				evt,
				"Poción activada. La XP de tu siguiente actividad que otorgue experiencia será el doble.",
			)
		} else {
			sendEcoSuccess(
				client,
				evt,
				"Escudo activado. Bloqueará el próximo intento de robo contra tu cuenta.",
			)
		}

	case "luck_charm":
		sendEcoSuccess(
			client,
			evt,
			"El Amuleto de Suerte se activa automáticamente al usar /slots. No hace falta consumirlo con /use.",
		)

	case "wood_chest", "gold_chest", "diamond_chest":
		sendEcoSuccess(
			client,
			evt,
			"Ese cofre se abre con /chest. No se ha consumido.",
		)

	case "pokeball", "master_ball":
		sendEcoError(
			client,
			evt,
			"El módulo de captura actual no está conectado a estos objetos. No se consumió nada.",
		)

	default:
		item, ok := store.GetItem(itemID)
		if !ok {
			sendEcoError(client, evt, "Ese objeto no tiene un efecto utilizable.")
			return
		}

		sendEcoError(
			client,
			evt,
			fmt.Sprintf("No hay un efecto de uso configurado para %s.", item.Name),
		)
	}
}