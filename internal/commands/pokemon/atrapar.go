package pokemon

import (
"fmt"
"math/rand"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/types/events"

"cometa-bot/internal/registry"
"cometa-bot/internal/store"
)

func init() {
registry.Register(&registry.Command{
Name:        "atrapar",
Aliases:     []string{"catch", "capturar"},
Description: "intenta atrapar el pokemon que apareció",
Category:    "𝐏⵿𝐨⵿𝐤⵿𝐞⵿𝐦⵿𝐨⵿𝐧",
Handler:     Atrapar,
})
}

func Atrapar(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()

spawn := store.GetSpawn(jid)
if spawn == nil {
sendText(client, evt, "❌ No hay ningún Pokémon para atrapar.\nUsa */poke* para hacer aparecer uno.")
return
}

if spawn.Cost > 0 {
balance := 0
if store.GlobalEconomy != nil {
balance = store.GlobalEconomy.GetBalance(jid)
}
if balance < spawn.Cost {
sendText(client, evt, fmt.Sprintf("❌ No tienes suficientes cometas.\nNecesitas *%d* pero solo tienes *%d*.\n\nTrabaja con */w* para ganar más.", spawn.Cost, balance))
return
}
store.GlobalEconomy.AddBalance(jid, -spawn.Cost)
}

success := rand.Intn(100) < 50

if success {
spawn.Pokemon.Date = "hoy"
store.AddPokemon(jid, spawn.Pokemon)
store.ClearSpawn(jid)

caption := fmt.Sprintf("╭〔 ✅ 𝐂𝐀𝐏𝐓𝐔𝐑𝐀𝐃𝐎 〕⬣\n┃\n┃ ¡Atrapaste a *%s*!\n┃ 🔢 *Nº:* %03d\n┃ ⚡ *Tipo:* %s\n┃\n┃ 🎉 Ya es parte de tu equipo\n┃\n╰━━━━━━━━━━━━━━━━⬣\n\n> Usa */pokeinv* para ver tu colección ☄️",
spawn.Pokemon.Name, spawn.Pokemon.ID, spawn.Pokemon.Type)

sendText(client, evt, caption)
} else {
store.ConsumeAttempt(jid)
spawn = store.GetSpawn(jid)

if spawn == nil {
sendText(client, evt, "❌ *¡Se escapó!*\n\nSe te acabaron los 3 intentos.\nUsa */poke* para buscar otro.")
} else {
costText := ""
if spawn.Cost > 0 {
costText = fmt.Sprintf("\n┃ 💰 Costará *%d* cometas", spawn.Cost)
}
sendText(client, evt, fmt.Sprintf("❌ *¡Se escapó!*\n\n┃ Intentos restantes: *%d/3*%s\n┃\nUsa */atrapar* de nuevo.", spawn.Attempts, costText))
}
}
}
