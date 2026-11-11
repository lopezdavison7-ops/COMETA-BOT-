package pokemon

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
Name:        "pokeinv",
Aliases:     []string{"mypokemon", "pokepokedex"},
Description: "muestra tu colección de pokemon atrapados",
Category:    "𝐏⵿𝐨⵿𝐤⵿𝐞⵿𝐦⵿𝐨⵿𝐧",
Handler:     PokeInv,
})
}

func PokeInv(evt *events.Message, client *whatsmeow.Client, _ []string) {
jid := evt.Info.Sender.String()
pokes := store.GetPokemonInventory(jid)

if len(pokes) == 0 {
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("╭〔 🎒 𝐏𝐎𝐊𝐄 𝐈𝐍𝐕 〕⬣\n\n┃ Aún no has atrapado ningún Pokémon.\n┃ ¡Usa */poke* para lanzar tu primera Pokéball!\n┃\n╰━━━━━━━━━━━━━━━━⬣"),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
return
}

b := strings.Builder{}
b.WriteString("╭〔 🎒 𝐏𝐎𝐊𝐄 𝐈𝐍𝐕 〕⬣\n┃\n")

// Mostrar máximo 15 para no saturar
maxShow := 15
for i, p := range pokes {
if i >= maxShow {
break
}
shiny := ""
if p.IsShiny {
shiny = " ✨"
}
b.WriteString(fmt.Sprintf("┃ *%d.* %s%s (Nº %03d)\n", i+1, p.Name, shiny, p.ID))
}

if len(pokes) > maxShow {
b.WriteString(fmt.Sprintf("┃ ...y %d Pokémon más\n", len(pokes)-maxShow))
}

b.WriteString(fmt.Sprintf("┃\n┃ 👥 Total atrapados: *%d*\n", len(pokes)))
b.WriteString("╰━━━━━━━━━━━━━━━━⬣")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
