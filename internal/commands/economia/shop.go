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

func init() {
registry.Register(&registry.Command{
Name:        "shop",
Aliases:     []string{"tienda", "store"},
Description: "muestra la tienda de items",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Shop,
})
registry.Register(&registry.Command{
Name:        "buy",
Aliases:     []string{"comprar", "purchase"},
Description: "compra un item de la tienda",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Buy,
})
registry.Register(&registry.Command{
Name:        "inv",
Aliases:     []string{"inventario", "items"},
Description: "muestra tus items",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Inv,
})
}

func Shop(evt *events.Message, client *whatsmeow.Client, _ []string) {
b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🛒 TIENDA` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

for i, item := range store.ShopItems {
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ *%d.* %s %s\n", i+1, item.Emoji, item.Name))
b.WriteString(fmt.Sprintf("   ↳ %s\n", item.Description))
b.WriteString(fmt.Sprintf("   ↳ 💰 *%d* cometas\n\n", item.Price))
}

b.WriteString("💡 Usá */buy <número>* para comprar\n")
b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}

func Buy(evt *events.Message, client *whatsmeow.Client, args []string) {
if len(args) == 0 {
sendEcoError(client, evt, "Usá: /buy <número del item>")
return
}

idx := 0
fmt.Sscanf(args[0], "%d", &idx)
if idx < 1 || idx > len(store.ShopItems) {
sendEcoError(client, evt, "Número de item inválido. Usá /shop para ver la lista.")
return
}

item := store.ShopItems[idx-1]
jid := evt.Info.Sender.String()
balance := store.GlobalEconomy.GetBalance(jid)

if balance < item.Price {
sendEcoError(client, evt, fmt.Sprintf("No tenés suficientes cometas.\nNecesitás *%d* pero tenés *%d*.", item.Price, balance))
return
}

store.GlobalEconomy.AddBalance(jid, -item.Price)
store.BuyItem(jid, item.ID)

sendEcoSuccess(client, evt, fmt.Sprintf("Compraste %s *%s* por %d cometas.\n\nUsá */inv* para ver tus items.", item.Emoji, item.Name, item.Price))
}

func Inv(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()
items := store.GetShopInventory(jid)

if len(items) == 0 {
sendEcoError(client, evt, "No tenés items.\nUsá /shop para comprar.")
return
}

counts := make(map[string]int)
for _, ui := range items {
counts[ui.ItemID]++
}

b := strings.Builder{}
b.WriteString("╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `🎒 INVENTARIO` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞\n> ㅤ\n")

for itemID, count := range counts {
item, ok := store.GetItem(itemID)
if !ok {
continue
}
b.WriteString(fmt.Sprintf("🫆֟፝𑑓᮫̤٘𑑓̸̷ %s *%s* x%d\n", item.Emoji, item.Name, count))
}

b.WriteString("╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
