package economia

import (
	"fmt"
	"strings"

	"cometa-bot/internal/store"
)

type economyEquipmentRequirement struct {
	Action      string
	ItemID      string
	Name        string
	Emoji       string
	Description string
	Price       int
}

var economyEquipment = []economyEquipmentRequirement{
	{
		Action:      "mine",
		ItemID:      "pickaxe",
		Name:        "Pico minero",
		Emoji:       "⛏️",
		Description: "Herramienta permanente necesaria para minar.",
		Price:       900,
	},
	{
		Action:      "fish",
		ItemID:      "fishing_rod",
		Name:        "Caña de pescar",
		Emoji:       "🎣",
		Description: "Herramienta permanente necesaria para pescar.",
		Price:       700,
	},
	{
		Action:      "hunt",
		ItemID:      "hunting_bow",
		Name:        "Arco de caza",
		Emoji:       "🏹",
		Description: "Herramienta permanente necesaria para cazar.",
		Price:       1100,
	},
	{
		Action:      "dig",
		ItemID:      "shovel",
		Name:        "Pala reforzada",
		Emoji:       "🪏",
		Description: "Herramienta permanente necesaria para excavar.",
		Price:       650,
	},
	{
		Action:      "talar",
		ItemID:      "axe",
		Name:        "Hacha de leñador",
		Emoji:       "🪓",
		Description: "Herramienta permanente necesaria para talar.",
		Price:       850,
	},
}

func init() {
	for _, equipment := range economyEquipment {
		store.ShopItems = append(store.ShopItems, store.Item{
			ID:          equipment.ItemID,
			Name:        equipment.Name,
			Description: equipment.Description,
			Price:       equipment.Price,
			Emoji:       equipment.Emoji,
		})
	}
}

func missingEconomyEquipmentMessage(jid, action string) string {
	for _, equipment := range economyEquipment {
		if equipment.Action != action {
			continue
		}

		if store.HasItem(jid, equipment.ItemID) {
			return ""
		}

		return fmt.Sprintf(
			"Necesitas %s para ejecutar /%s.\n"+
				"Consíguelo en /shop y cómpralo con /buy <número>.\n"+
				"La herramienta permanece en tu inventario y no se consume.",
			equipment.Name,
			action,
		)
	}

	return ""
}

func missingAllwEquipmentMessage(jid string) string {
	var missing []string

	for _, equipment := range economyEquipment {
		if !store.HasItem(jid, equipment.ItemID) {
			missing = append(
				missing,
				equipment.Emoji+" "+equipment.Name,
			)
		}
	}

	if len(missing) == 0 {
		return ""
	}

	return "Para ejecutar /allw necesitas todas las herramientas de sus actividades:\n"+
		strings.Join(missing, "\n")+
		"\nConsigue las que te faltan en /shop y cómpralas con /buy <número>."
}