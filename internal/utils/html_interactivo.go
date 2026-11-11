package utils

import (
"crypto/rand"
"encoding/json"
"fmt"

"go.mau.fi/whatsmeow/proto/waE2E"
"google.golang.org/protobuf/reflect/protoreflect"
)

func randomUUID() string {
b := make([]byte, 16)
rand.Read(b)
b[6] = (b[6] & 0x0f) | 0x40
b[8] = (b[8] & 0x3f) | 0x80
return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func field(m protoreflect.Message, name string) protoreflect.FieldDescriptor {
return m.Descriptor().Fields().ByName(protoreflect.Name(name))
}

func setNumber(m protoreflect.Message, name string, v int32) {
fd := field(m, name)
if fd == nil { return }
switch fd.Kind() {
case protoreflect.EnumKind:
m.Set(fd, protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)))
case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
m.Set(fd, protoreflect.ValueOfInt32(v))
case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
m.Set(fd, protoreflect.ValueOfInt64(int64(v)))
case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
m.Set(fd, protoreflect.ValueOfUint32(uint32(v)))
}
}

func setBool(m protoreflect.Message, name string, v bool) {
fd := field(m, name)
if fd == nil { return }
if fd.Kind() == protoreflect.BoolKind {
m.Set(fd, protoreflect.ValueOfBool(v))
}
}

func BuildUnifiedResponseJSON(htmlPayload string) string {
payload := map[string]interface{}{
"__typename":  "GenAIUnifiedResponse",
"response_id": randomUUID(),
"sections": []interface{}{
map[string]interface{}{
"__typename": "GenAIUnifiedResponseSection",
"view_model": map[string]interface{}{
"__typename": "GenAISingleLayoutViewModel",
"primitive": map[string]interface{}{
"__typename":      "FOAHtmlPrimitiveDemoDONOTUSE",
"trusted_sources": []string{},
"payload":         htmlPayload,
},
},
},
},
}
b, _ := json.Marshal(payload)
return string(b)
}

func BuildBotForwardedMessage(htmlPayload string) *waE2E.Message {
jsonStr := BuildUnifiedResponseJSON(htmlPayload)

outer := &waE2E.Message{}
or := outer.ProtoReflect()

fdBot := field(or, "botForwardedMessage")
if fdBot == nil || fdBot.Message() == nil { return nil }
botMsg := or.NewField(fdBot).Message()

fdInner := field(botMsg, "message")
if fdInner == nil || fdInner.Message() == nil { return nil }
inner := botMsg.NewField(fdInner).Message()

fdRich := field(inner, "richResponseMessage")
if fdRich == nil || fdRich.Message() == nil { return nil }
rich := inner.NewField(fdRich).Message()

setNumber(rich, "messageType", 1)

if fdUR := field(rich, "unifiedResponse"); fdUR != nil {
if fdUR.Message() != nil {
ur := rich.NewField(fdUR).Message()
if fdData := field(ur, "data"); fdData != nil {
switch fdData.Kind() {
case protoreflect.BytesKind:
ur.Set(fdData, protoreflect.ValueOfBytes([]byte(jsonStr)))
case protoreflect.StringKind:
ur.Set(fdData, protoreflect.ValueOfString(jsonStr))
}
rich.Set(fdUR, protoreflect.ValueOfMessage(ur))
}
} else if fdUR.Kind() == protoreflect.BytesKind {
rich.Set(fdUR, protoreflect.ValueOfBytes([]byte(jsonStr)))
}
}

if fdCI := field(rich, "contextInfo"); fdCI != nil && fdCI.Message() != nil {
ci := rich.NewField(fdCI).Message()
setBool(ci, "isForwarded", true)
setNumber(ci, "forwardOrigin", 4)
rich.Set(fdCI, protoreflect.ValueOfMessage(ci))
}

inner.Set(fdRich, protoreflect.ValueOfMessage(rich))
botMsg.Set(fdInner, protoreflect.ValueOfMessage(inner))
or.Set(fdBot, protoreflect.ValueOfMessage(botMsg))

return outer
}
