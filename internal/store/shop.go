package store

import (
"encoding/json"
	"fmt"
"os"
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
ShopItems   = []Item{
{ID: "luck_charm", Name: "Amuleto de Suerte", Description: "+10% prob en slots", Price: 500, Emoji: "🍀"},
{ID: "shield", Name: "Escudo Protector", Description: "Protege de 1 robo", Price: 800, Emoji: "🛡️"},
{ID: "xp_boost", Name: "Poción de XP", Description: "Doble XP por 1 trabajo", Price: 300, Emoji: "⚗️"},
{ID: "wood_chest", Name: "Cofre de Madera", Description: "Recompensa aleatoria", Price: 200, Emoji: "📦"},
{ID: "gold_chest", Name: "Cofre de Oro", Description: "Recompensa épica", Price: 1000, Emoji: "🎁"},
{ID: "diamond_chest", Name: "Cofre Diamante", Description: "Recompensa legendaria", Price: 3000, Emoji: "💎"},
{ID: "pokeball", Name: "Pokéball Extra", Description: "+1 intento para atrapar", Price: 400, Emoji: "🔴"},
{ID: "master_ball", Name: "Master Ball", Description: "Captura garantizada", Price: 5000, Emoji: "🟣"},
}

UserInventories = make(map[string][]UserItem)
shopMu          sync.Mutex
shopFile        = "database/shop.json"
)

func init() {
os.MkdirAll("database", 0755)
loadShop()
}

func loadShop() {
shopMu.Lock()
defer shopMu.Unlock()
data, err := os.ReadFile(shopFile)
if err == nil {
json.Unmarshal(data, &UserInventories)
}
}

func saveShop() {
shopMu.Lock()
defer shopMu.Unlock()
data, _ := json.MarshalIndent(UserInventories, "", "  ")
os.WriteFile(shopFile, data, 0644)
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
UserInventories[jid] = append(UserInventories[jid], UserItem{ItemID: itemID, Amount: 1})
shopMu.Unlock()
saveShop()
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
if ui.ItemID == itemID {
return true
}
}
return false
}

func UseItem(jid, itemID string) bool {
shopMu.Lock()
defer shopMu.Unlock()
for i, ui := range UserInventories[jid] {
if ui.ItemID == itemID {
UserInventories[jid] = append(UserInventories[jid][:i], UserInventories[jid][i+1:]...)
saveShop()
return true
}
}
return false
}
