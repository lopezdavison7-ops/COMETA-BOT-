package economia

import (
	"strconv"
	"strings"

	"cometa-bot/internal/store"
)

func resolveInventoryItem(jid, selector string) string {
	selector = strings.TrimSpace(strings.ToLower(selector))

	if selector == "" {
		return ""
	}

	if _, ok := store.GetItem(selector); ok && store.HasItem(jid, selector) {
		return selector
	}

	if index, err := strconv.Atoi(selector); err == nil && index > 0 {
		position := 0

		for _, inventoryItem := range store.GetShopInventory(jid) {
			for count := 0; count < inventoryItem.Amount; count++ {
				position++

				if position == index {
					return inventoryItem.ItemID
				}
			}
		}
	}

	for _, inventoryItem := range store.GetShopInventory(jid) {
		item, ok := store.GetItem(inventoryItem.ItemID)
		if ok && (
			strings.ToLower(item.Name) == selector ||
				strings.ToLower(item.ID) == selector,
		) {
			return inventoryItem.ItemID
		}
	}

	return ""
}