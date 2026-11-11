package handler

import (
"context"
"fmt"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types"
"google.golang.org/protobuf/proto"
)

var ownerJIDs = []string{
"50588910385@s.whatsapp.net",
"50578391933@s.whatsapp.net",
"50577455383@s.whatsapp.net",
}

func NotifyOwners(client *whatsmeow.Client, msg string) {
if client == nil || !client.IsConnected() {
fmt.Println("⚠️ NotifyOwners: cliente no disponible")
return
}
for _, jidStr := range ownerJIDs {
jid, err := types.ParseJID(jidStr)
if err != nil {
continue
}
client.SendMessage(context.Background(), jid, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text: proto.String("🔔 Notificación del sistema:\n\n" + msg),
},
})
}
}
