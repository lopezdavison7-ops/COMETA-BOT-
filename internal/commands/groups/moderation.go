package groups

import (
"context"
"fmt"
"regexp"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "kick",
Aliases:     []string{"sacar", "ban"},
Description: "saca un miembro del grupo",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Kick,
})
registry.Register(&registry.Command{
Name:        "add",
Aliases:     []string{"agregar", "invitar"},
Description: "agrega un miembro al grupo",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Add,
})
registry.Register(&registry.Command{
Name:        "promote",
Aliases:     []string{"daradmin", "promover"},
Description: "da admin a un miembro",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Promote,
})
registry.Register(&registry.Command{
Name:        "demote",
Aliases:     []string{"quitaradmin"},
Description: "quita admin a un miembro",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Demote,
})
registry.Register(&registry.Command{
Name:        "tagall",
Aliases:     []string{"todos", "all"},
Description: "menciona a todos los miembros",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     TagAll,
})
registry.Register(&registry.Command{
Name:        "hidetag",
Aliases:     []string{"htag"},
Description: "mención invisible a todos",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     Hidetag,
})
registry.Register(&registry.Command{
Name:        "grouplink",
Aliases:     []string{"link", "linkgc"},
Description: "obtiene el link del grupo",
Category:    "𝐆⵿𝐫⵿𝐮⵿𝐩⵿𝐨⵿𝐬",
Handler:     GroupLink,
})
}

func getTargetFromMsg(evt *events.Message, args []string) string {
ext := evt.Message.GetExtendedTextMessage()
if ext != nil {
ctx := ext.GetContextInfo()
if ctx != nil {
mentioned := ctx.GetMentionedJID()
if len(mentioned) > 0 {
return mentioned[0]
}
if p := ctx.GetParticipant(); p != "" {
return p
}
}
}
if len(args) > 0 {
return args[0]
}
return ""
}

func cleanNumber(raw string) string {
re := regexp.MustCompile(`[^\d]`)
return re.ReplaceAllString(raw, "")
}

func Kick(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
target := getTargetFromMsg(evt, args)
if target == "" {
sendGroupError(client, evt, "Mencioná a alguien o respondé a su mensaje.")
return
}
jid, err := types.ParseJID(target)
if err != nil {
sendGroupError(client, evt, "Usuario inválido.")
return
}
_, err = client.UpdateGroupParticipants(context.Background(), evt.Info.Chat, []types.JID{jid}, "remove")
if err != nil {
sendGroupError(client, evt, "Error al sacar: "+err.Error())
return
}
sendGroupSuccess(client, evt, "Usuario sacado del grupo.")
}

func Add(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
if len(args) == 0 {
sendGroupError(client, evt, "Usá: /add <número>\nEjemplo: /add +503 7563 8328")
return
}

// Unir todos los argumentos y limpiar caracteres no numéricos
raw := strings.Join(args, " ")
num := cleanNumber(raw)

if len(num) < 7 {
sendGroupError(client, evt, "Número inválido. Debe tener al menos 7 dígitos.")
return
}

jid := types.NewJID(num, types.DefaultUserServer)
_, err := client.UpdateGroupParticipants(context.Background(), evt.Info.Chat, []types.JID{jid}, "add")
if err != nil {
sendGroupError(client, evt, "Error: "+err.Error())
return
}
sendGroupSuccess(client, evt, "Usuario agregado al grupo.")
}

func Promote(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
target := getTargetFromMsg(evt, args)
if target == "" {
sendGroupError(client, evt, "Mencioná a alguien.")
return
}
jid, err := types.ParseJID(target)
if err != nil {
sendGroupError(client, evt, "Usuario inválido.")
return
}
_, err = client.UpdateGroupParticipants(context.Background(), evt.Info.Chat, []types.JID{jid}, "promote")
if err != nil {
sendGroupError(client, evt, "Error: "+err.Error())
return
}
sendGroupSuccess(client, evt, "Usuario promovido a admin.")
}

func Demote(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
target := getTargetFromMsg(evt, args)
if target == "" {
sendGroupError(client, evt, "Mencioná a alguien.")
return
}
jid, err := types.ParseJID(target)
if err != nil {
sendGroupError(client, evt, "Usuario inválido.")
return
}
_, err = client.UpdateGroupParticipants(context.Background(), evt.Info.Chat, []types.JID{jid}, "demote")
if err != nil {
sendGroupError(client, evt, "Error: "+err.Error())
return
}
sendGroupSuccess(client, evt, "Usuario degradado.")
}

func TagAll(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
group, err := client.GetGroupInfo(context.Background(), evt.Info.Chat)
if err != nil {
sendGroupError(client, evt, "Error al obtener info del grupo.")
return
}
msg := "╭〔 👥 𝐓𝐀𝐆𝐀𝐋𝐋 〕⬣\n"
if len(args) > 0 {
msg += "┃ 📢 " + strings.Join(args, " ") + "\n"
}
msg += "┃\n"
var mentions []string
for i, p := range group.Participants {
msg += fmt.Sprintf("┃ %d. @%s\n", i+1, p.JID.User)
mentions = append(mentions, p.JID.String())
}
msg += "╰━━━━━━━━━━━━━━━━⬣"
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(msg),
ContextInfo: &waE2E.ContextInfo{
MentionedJID:  mentions,
QuotedMessage: evt.Message,
},
},
})
}

func Hidetag(evt *events.Message, client *whatsmeow.Client, args []string) {
if !requireAdmin(client, evt) {
return
}
group, err := client.GetGroupInfo(context.Background(), evt.Info.Chat)
if err != nil {
sendGroupError(client, evt, "Error al obtener info del grupo.")
return
}
text := strings.Join(args, " ")
if text == "" {
text = "📢"
}
var mentions []string
for _, p := range group.Participants {
mentions = append(mentions, p.JID.String())
}
client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String(text),
ContextInfo: &waE2E.ContextInfo{
MentionedJID:  mentions,
QuotedMessage: evt.Message,
},
},
})
}

func GroupLink(evt *events.Message, client *whatsmeow.Client, _ []string) {
if !requireAdmin(client, evt) {
return
}
link, err := client.GetGroupInviteLink(context.Background(), evt.Info.Chat, false)
if err != nil {
sendGroupError(client, evt, "Error al obtener link.")
return
}
sendGroupMsg(client, evt.Info.Chat, "╭〔 🔗 𝐋𝐈𝐍𝐊 𝐃𝐄𝐋 𝐆𝐑𝐔𝐏𝐎 〕⬣\n┃\n┃ "+link+"\n┃\n╰━━━━━━━━━━━━━━━━⬣", evt.Message)
}
