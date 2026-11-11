package system

import (
"context"
"fmt"
"strings"

"go.mau.fi/whatsmeow"
"go.mau.fi/whatsmeow/proto/waE2E"
"go.mau.fi/whatsmeow/types/events"
"google.golang.org/protobuf/proto"

"cometa-bot/internal/registry"
)

func init() {
registry.Register(&registry.Command{
Name:        "ext",
Aliases:     []string{"extract", "metadata"},
Description: "extrae los metadatos de un mensaje (citado o actual)",
Category:    "𝐒⵿𝐢⵿𝐬⵿𝐭⵿𝐞⵿𝐦⵿𝐚",
Handler:     Extractor,
})
}

func extraerContextInfo(msg *waE2E.Message) *waE2E.ContextInfo {
if msg == nil {
return nil
}
if ext := msg.GetExtendedTextMessage(); ext != nil {
return ext.GetContextInfo()
}
if img := msg.GetImageMessage(); img != nil {
return img.GetContextInfo()
}
if vid := msg.GetVideoMessage(); vid != nil {
return vid.GetContextInfo()
}
if doc := msg.GetDocumentMessage(); doc != nil {
return doc.GetContextInfo()
}
if stk := msg.GetStickerMessage(); stk != nil {
return stk.GetContextInfo()
}
if aud := msg.GetAudioMessage(); aud != nil {
return aud.GetContextInfo()
}
return nil
}

func extraerTexto(msg *waE2E.Message) (string, string) {
if msg == nil {
return "desconocido", ""
}
switch {
case msg.GetConversation() != "":
return "texto", msg.GetConversation()
case msg.GetExtendedTextMessage() != nil:
return "texto (extendido)", msg.GetExtendedTextMessage().GetText()
case msg.GetImageMessage() != nil:
return "imagen", msg.GetImageMessage().GetCaption()
case msg.GetVideoMessage() != nil:
return "video", msg.GetVideoMessage().GetCaption()
case msg.GetDocumentMessage() != nil:
return "documento", msg.GetDocumentMessage().GetCaption()
case msg.GetStickerMessage() != nil:
return "sticker", ""
case msg.GetAudioMessage() != nil:
return "audio", ""
case msg.GetContactMessage() != nil:
return "contacto", msg.GetContactMessage().GetDisplayName()
case msg.GetLocationMessage() != nil:
return "ubicación", ""
default:
return "otro", ""
}
}

func Extractor(evt *events.Message, client *whatsmeow.Client, _ []string) {
// ContextInfo del mensaje ACTUAL (contiene la referencia al citado)
currentCtx := extraerContextInfo(evt.Message)

// Mensaje citado
quoted := currentCtx.GetQuotedMessage()
origen := "ACTUAL"
target := evt.Message

if quoted != nil {
origen = "CITADO"
target = quoted
}

msgType, textContent := extraerTexto(target)

// 🔑 CLAVE: El ID y remitente del citado están en el ContextInfo ACTUAL
msgID := "no disponible"
sender := ""

if origen == "CITADO" && currentCtx != nil {
if currentCtx.GetStanzaID() != "" {
msgID = currentCtx.GetStanzaID()
}
sender = currentCtx.GetParticipant()
} else if currentCtx != nil {
if currentCtx.GetStanzaID() != "" {
msgID = currentCtx.GetStanzaID()
}
sender = evt.Info.Sender.String()
}

var b strings.Builder
b.WriteString("╭─ ✦ ˚ ✦ ˚ ✦ ─╮\n\n")
b.WriteString("🔎 *COMETA-BOT · EXTRACTOR DE MENSAJES*\n\n")
b.WriteString("*CONTENIDO*\n")
b.WriteString(fmt.Sprintf("• Origen: `%s`\n", origen))
b.WriteString(fmt.Sprintf("• Tipo: `%s`\n", msgType))
b.WriteString(fmt.Sprintf("• ID del mensaje: `%s`\n", msgID))
b.WriteString(fmt.Sprintf("• Remitente: `%s`\n", sender))
if textContent != "" {
b.WriteString(fmt.Sprintf("• Texto/caption: %s\n", textContent))
}
b.WriteString("\n")

// Buscar metadatos de reenvío en el mensaje target
targetCtx := extraerContextInfo(target)
encontrado := false

if targetCtx != nil {
if targetCtx.GetIsForwarded() || targetCtx.GetForwardingScore() > 0 {
b.WriteString("*REENVÍO Y ATRIBUCIÓN*\n")
b.WriteString(fmt.Sprintf("• Marcado como reenviado: `%v`\n", targetCtx.GetIsForwarded()))
b.WriteString(fmt.Sprintf("• Puntuación de reenvío: `%d`\n", targetCtx.GetForwardingScore()))
encontrado = true
}

if nl := targetCtx.GetForwardedNewsletterMessageInfo(); nl != nil {
b.WriteString("\n*METADATOS DE CANAL*\n")
b.WriteString(fmt.Sprintf("• Newsletter JID: `%s`\n", nl.GetNewsletterJID()))
b.WriteString(fmt.Sprintf("• Newsletter Name: `%s`\n", nl.GetNewsletterName()))
b.WriteString(fmt.Sprintf("• Server Message ID: `%d`\n", nl.GetServerMessageID()))
encontrado = true
}
}

// Si no encontró en el target, buscar en el mensaje actual
if !encontrado && currentCtx != nil && origen == "CITADO" {
if currentCtx.GetIsForwarded() || currentCtx.GetForwardingScore() > 0 {
b.WriteString("*REENVÍO Y ATRIBUCIÓN (mensaje actual)*\n")
b.WriteString(fmt.Sprintf("• Marcado como reenviado: `%v`\n", currentCtx.GetIsForwarded()))
b.WriteString(fmt.Sprintf("• Puntuación de reenvío: `%d`\n", currentCtx.GetForwardingScore()))
encontrado = true
}

if nl := currentCtx.GetForwardedNewsletterMessageInfo(); nl != nil {
b.WriteString("\n*METADATOS DE CANAL (mensaje actual)*\n")
b.WriteString(fmt.Sprintf("• Newsletter JID: `%s`\n", nl.GetNewsletterJID()))
b.WriteString(fmt.Sprintf("• Newsletter Name: `%s`\n", nl.GetNewsletterName()))
b.WriteString(fmt.Sprintf("• Server Message ID: `%d`\n", nl.GetServerMessageID()))
encontrado = true
}
}

if !encontrado {
b.WriteString("*METADATOS*\n")
b.WriteString("• No se encontraron metadatos de reenvío ni de canal.\n")
}

b.WriteString("\n> _Nota: WhatsApp puede quitar metadatos cuando un contenido se copia, descarga o vuelve a enviar de otra forma._\n\n")
b.WriteString("╰─ ✦ ˚ ✦ ˚ ✦ ─╯")

client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
ExtendedTextMessage: &waE2E.ExtendedTextMessage{
Text:        proto.String(b.String()),
ContextInfo: &waE2E.ContextInfo{QuotedMessage: evt.Message},
},
})
}
