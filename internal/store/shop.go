package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Emoji       string `json:"emoji"`
}

type UserItem struct {
	ItemID string `json:"item_id"`
	Amount int    `json:"amount"`
}

var (
	ShopItems = []Item{
		{ID: "luck_charm", Name: "Amuleto de Suerte", Description: "+10% probabilidad en slots; se consume automáticamente", Price: 500, Emoji: "🍀"},
		{ID: "shield", Name: "Escudo Protector", Description: "Bloquea un intento de robo", Price: 800, Emoji: "🛡️"},
		{ID: "xp_boost", Name: "Poción de XP", Description: "Duplica la XP de tu siguiente actividad con XP", Price: 300, Emoji: "⚗️"},
		{ID: "wood_chest", Name: "Cofre de Madera", Description: "Recompensa aleatoria", Price: 200, Emoji: "📦"},
		{ID: "gold_chest", Name: "Cofre de Oro", Description: "Recompensa épica", Price: 1000, Emoji: "🎁"},
		{ID: "diamond_chest", Name: "Cofre Diamante", Description: "Recompensa legendaria", Price: 3000, Emoji: "💎"},
		{ID: "pokeball", Name: "Pokéball Extra", Description: "+1 intento para atrapar", Price: 400, Emoji: "🔴"},
		{ID: "master_ball", Name: "Master Ball", Description: "Captura garantizada", Price: 5000, Emoji: "🟣"},
	}
	UserInventories = make(map[string][]UserItem)
	shopMu          sync.Mutex
	shopFile        = filepath.Join("database", "shop.json")
)

func init() {
	_ = os.MkdirAll(filepath.Dir(shopFile), 0755)
	loadShop()
}

func loadShop() {
	shopMu.Lock()
	defer shopMu.Unlock()

	data, err := os.ReadFile(shopFile)
	if err != nil {
		return
	}

	var loaded map[string][]UserItem
	if err := json.Unmarshal(data, &loaded); err != nil {
		return
	}

	if loaded == nil {
		loaded = make(map[string][]UserItem)
	}

	UserInventories = loaded
}

func saveShopLocked() error {
	data, err := json.MarshalIndent(UserInventories, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(shopFile), 0755); err != nil {
		return err
	}

	tmp := shopFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmp, shopFile); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func GetItem(id string) (Item, bool) {
	for _, item := range ShopItems {
		if item.ID == id {
			return item, true
		}
	}

	return Item{}, false
}

func BuyItem(jid, itemID string) (Item, error) {
	item, ok := GetItem(itemID)
	if !ok {
		return Item{}, fmt.Errorf("item no existe")
	}

	shopMu.Lock()
	defer shopMu.Unlock()

	UserInventories[jid] = append(
		UserInventories[jid],
		UserItem{ItemID: itemID, Amount: 1},
	)

	if err := saveShopLocked(); err != nil {
		items := UserInventories[jid]
		if len(items) > 0 {
			UserInventories[jid] = items[:len(items)-1]
		}

		return Item{}, err
	}

	return item, nil
}

func GetShopInventory(jid string) []UserItem {
	shopMu.Lock()
	defer shopMu.Unlock()

	res := make([]UserItem, len(UserInventories[jid]))
	copy(res, UserInventories[jid])

	return res
}

func HasItem(jid, itemID string) bool {
	shopMu.Lock()
	defer shopMu.Unlock()

	for _, ui := range UserInventories[jid] {
		if ui.ItemID == itemID && ui.Amount > 0 {
			return true
		}
	}

	return false
}

func removeOneItemLocked(jid, itemID string) bool {
	items := UserInventories[jid]

	for i, ui := range items {
		if ui.ItemID == itemID && ui.Amount > 0 {
			if ui.Amount > 1 {
				UserInventories[jid][i].Amount--
			} else {
				UserInventories[jid] = append(items[:i], items[i+1:]...)
			}

			return true
		}
	}

	return false
}

func UseItem(jid, itemID string) bool {
	shopMu.Lock()
	defer shopMu.Unlock()

	before := append([]UserItem(nil), UserInventories[jid]...)

	if !removeOneItemLocked(jid, itemID) {
		return false
	}

	if err := saveShopLocked(); err != nil {
		UserInventories[jid] = before
		return false
	}

	return true
}

func SellItem(jid, itemID string) (Item, error) {
	item, ok := GetItem(itemID)
	if !ok {
		return Item{}, fmt.Errorf("item no existe")
	}

	shopMu.Lock()
	defer shopMu.Unlock()

	before := append([]UserItem(nil), UserInventories[jid]...)

	if !removeOneItemLocked(jid, itemID) {
		return Item{}, fmt.Errorf("no tienes ese objeto")
	}

	if err := saveShopLocked(); err != nil {
		UserInventories[jid] = before
		return Item{}, err
	}

	return item, nil
}