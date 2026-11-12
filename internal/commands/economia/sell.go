package economia

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	"cometa-bot/internal/registry"
	"cometa-bot/internal/store"
)

func init() {
	registry.Register(&registry.Command{
		Name:        "sellitem",
		Aliases:     []string{"venderobjeto", "sellobjeto"},
		Description: "vende un objeto de tu inventario por la mitad de su precio",
		Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
		Handler:     Sell,
	})
}

func Sell(evt *events.Message, client *whatsmeow.Client, args []string) {
	if store.GlobalEconomy == nil {
		sendEcoError(client, evt, "La economía no está inicializada.")
		return
	}

	if len(args) == 0 {
		sendEcoError(
			client,
			evt,
			"Usá: /sellitem <id del objeto>. El valor de venta es el 50% del precio de tienda.",
		)
		return
	}

	jid := evt.Info.Sender.String()
	itemID := resolveInventoryItem(jid, args[0])

	if itemID == "" {
		sendEcoError(client, evt, "No encontré ese objeto en tu inventario. Revisá /ecoinv.")
		return
	}

	item, err := store.SellItem(jid, itemID)
	if err != nil {
		sendEcoError(client, evt, "No se pudo vender el objeto: "+err.Error())
		return
	}

	price := item.Price / 2
	if price < 1 {
		price = 1
	}

	store.GlobalEconomy.AddBalance(jid, price)

	sendEcoSuccess(
		client,
		evt,
		fmt.Sprintf(
			"Vendiste %s *%s* por *%d* cometas. El objeto fue retirado del inventario.",
			item.Emoji,
			item.Name,
			price,
		),
	)
}