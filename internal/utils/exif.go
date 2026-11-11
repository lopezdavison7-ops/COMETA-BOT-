package utils

import (
"bytes"
"crypto/rand"
"encoding/binary"
"encoding/hex"
"encoding/json"
)

func validWebp(buf []byte) bool {
if len(buf) < 12 || string(buf[0:4]) != "RIFF" || string(buf[8:12]) != "WEBP" {
return false
}
total := int(binary.LittleEndian.Uint32(buf[4:8])) + 8
if total != len(buf) {
return false
}
off := 12
for off < len(buf) {
if off+8 > len(buf) {
return false
}
size := int(binary.LittleEndian.Uint32(buf[off+4 : off+8]))
off += 8 + size
if size%2 == 1 {
off++
}
}
return off == len(buf)
}

func AddExif(webpBuffer []byte, packname, author string) ([]byte, error) {
if !validWebp(webpBuffer) {
return webpBuffer, nil
}

packID := make([]byte, 32)
rand.Read(packID)

jsonData := map[string]interface{}{
"sticker-pack-id":        hex.EncodeToString(packID),
"sticker-pack-name":      packname,
"sticker-pack-publisher": author,
"emojis":                 []string{"☄️"},
}

jsonBytes, err := json.Marshal(jsonData)
if err != nil {
return webpBuffer, nil
}

exifAttr := []byte{
0x49, 0x49, 0x2A, 0x00,
0x08, 0x00, 0x00, 0x00,
0x01, 0x00, 0x41, 0x57,
0x07, 0x00, 0x00, 0x00,
0x00, 0x00, 0x16, 0x00,
0x00, 0x00,
}

exif := append(append([]byte{}, exifAttr...), jsonBytes...)
binary.LittleEndian.PutUint32(exif[14:18], uint32(len(jsonBytes)))

var chunk bytes.Buffer
chunk.WriteString("EXIF")
sizeBuf := make([]byte, 4)
binary.LittleEndian.PutUint32(sizeBuf, uint32(len(exif)))
chunk.Write(sizeBuf)
chunk.Write(exif)
if len(exif)%2 == 1 {
chunk.WriteByte(0)
}

body := append(append([]byte{}, webpBuffer[12:]...), chunk.Bytes()...)

var out bytes.Buffer
out.WriteString("RIFF")
totalSize := 4 + len(body)
ts := make([]byte, 4)
binary.LittleEndian.PutUint32(ts, uint32(totalSize))
out.Write(ts)
out.WriteString("WEBP")
out.Write(body)

result := out.Bytes()
if !validWebp(result) {
return webpBuffer, nil
}

return result, nil
}
