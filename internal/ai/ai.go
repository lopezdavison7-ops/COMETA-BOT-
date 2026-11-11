package ai

import (
"bytes"
"crypto/md5"
"crypto/rand"
"encoding/hex"
"encoding/json"
"fmt"
"io"
"mime/multipart"
"net/http"
"strings"
"time"
)

var httpClient = &http.Client{Timeout: 20 * time.Second}

func reverseString(s string) string {
runes := []rune(s)
for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
runes[i], runes[j] = runes[j], runes[i]
}
return string(runes)
}

func md5Reverse(input string) string {
hash := md5.Sum([]byte(input))
hexStr := hex.EncodeToString(hash[:])
return reverseString(hexStr)
}

func generateIslandKey(userAgent string) string {
b := make([]byte, 8)
rand.Read(b)
myrandomstr := fmt.Sprintf("%d", int64(b[0])<<56|int64(b[1])<<48|int64(b[2])<<40|int64(b[3])<<32|int64(b[4])<<24|int64(b[5])<<16|int64(b[6])<<8|int64(b[7]))
if len(myrandomstr) > 11 {
myrandomstr = myrandomstr[len(myrandomstr)-11:]
}

inner := md5Reverse(userAgent + myrandomstr + "hackers_become_a_little_stinkier_every_time_they_hack")
middle := md5Reverse(userAgent + inner)
outer := md5Reverse(userAgent + middle)

return "tryit-" + myrandomstr + "-" + outer
}

func generateUUID() string {
uuid := make([]byte, 16)
rand.Read(uuid)
uuid[6] = (uuid[6] & 0x0f) | 0x40
uuid[8] = (uuid[8] & 0x3f) | 0x80
return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func DeepAIChat(prompt string) (string, error) {
userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
key := generateIslandKey(userAgent)
sessionUUID := generateUUID()
sensitivityRequestID := generateUUID()

messages := []map[string]string{
{"role": "system", "content": "Eres Cometa Bot, asistente de Nicaragua. Responde corto y con jerga nica."},
{"role": "user", "content": prompt},
}

messagesJSON, _ := json.Marshal(messages)

var b bytes.Buffer
w := multipart.NewWriter(&b)

w.WriteField("chat_style", "chat")
w.WriteField("model", "standard")
w.WriteField("session_uuid", sessionUUID)
w.WriteField("sensitivity_request_id", sensitivityRequestID)
w.WriteField("hacker_is_stinky", "very_stinky")
w.WriteField("enabled_tools", `["image_generator", "image_editor"]`)
w.WriteField("chatHistory", string(messagesJSON))

w.Close()

req, _ := http.NewRequest("POST", "https://api.deepai.org/hacking_is_a_serious_crime", &b)
req.Header.Set("Content-Type", w.FormDataContentType())
req.Header.Set("api-key", key)
req.Header.Set("user-agent", userAgent)
req.Header.Set("referer", "https://deepai.org/chat")
req.Header.Set("origin", "https://deepai.org")

resp, err := httpClient.Do(req)
if err != nil {
return "", err
}
defer resp.Body.Close()

if resp.StatusCode != 200 {
return "", fmt.Errorf("HTTP %d", resp.StatusCode)
}

body, _ := io.ReadAll(resp.Body)
return strings.TrimSpace(string(body)), nil
}

func Ask(prompt string) string {
res, err := DeepAIChat(prompt)
if err != nil || res == "" {
return "Uff chacho, ando un poco saturado. Manda de nuevo en unos segundos ️"
}
return res
}
