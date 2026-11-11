package economia

import (
"context"
"fmt"
"math/rand"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "chest",
Aliases:     []string{"cofre", "abrir"},
Description: "abre un cofre que tengas en tu inventario",
Category:    "𝐄⵿𝐜𝵿𝐨⵿𝐧𝐨⵿𝐦⵿𝐢⵿𝐚",
Handler:     Chest,
})
}

func Chest(evt *events.Message, client *whatsmeow.Client, args []string) {
jid := evt.Info.Sender.String()

chestTypes := []string{"wood_chest", "gold_chest", "diamond_chest"}
var usedChest string

for _, ct := range chestTypes {
if store.HasItem(jid, ct) {
usedChest = ct
break
}
}

if usedChest == "" {
sendEcoError(client, evt, "No tenés cofres.\nComprá uno en /shop")
return
}

store.UseItem(jid, usedChest)

var minReward, maxReward int
var chestName, chestEmoji string

switch usedChest {
case "wood_chest":
minReward, maxReward = 50, 200
chestName = "Cofre de Madera"
chestEmoji = "📦"
case "gold_chest":
minReward, maxReward = 300, 800
chestName = "Cofre de Oro"
chestEmoji = "🎁"
case "diamond_chest":
minReward, maxReward = 1000, 3000
chestName = "Cofre Diamante"
chestEmoji = "💎"
}

reward := rand.Intn(maxReward-minReward+1) + minReward
store.GlobalEconomy.AddBalance(jid, reward)

bonus := ""
if rand.Intn(100) < 20 {
bonusXP := rand.Intn(50) + 20
bonus = fmt.Sprintf("\n🫆֟፝𑑓᮫̤٘𑑓̸̷ *Bonus:* +%d XP", bonusXP)
}

text := fmt.Sprintf(`╭ᩙ╌͜╍͜╌͜─͜─͜─͜─፝͜「 `+"`%s COFRE`"+` 」፝͜─͜─͜─͜─͜╌͜━㍃゙⃞
> ㅤ
🫆֟፝𑑓᮫̤٘𑑓̸̷ Abriste un *%s*
🫆֟፝𑑓᮫̤٘𑑓̸̷ 💰 *Recompensa:* %d cometas%s
╰᳕╌͜╍͜╌͜─͜─͜─͜─͜─͜─⵿֟፝─͜─͜─͜─͜─͜─͜─͜─͜╌͜━㍃゙⃞`, chestEmoji, chestName, reward, bonus)

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
