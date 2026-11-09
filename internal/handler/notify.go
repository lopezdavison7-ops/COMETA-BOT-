package handler

import (
"context"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"google.golang.org/protobuf/proto"
)

var Owners = []string{"50577455383", "50578391933"}

func NotifyOwners(client *whatsmeow.Client, text string) {
for _, num := range Owners {
jid := types.NewJID(num, types.DefaultUserServer)
client.SendMessage(context.Background(), jid, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("🚨 *AUTO-REPAIR / ALERTA*\n\n" + text),
},
})
}
}
