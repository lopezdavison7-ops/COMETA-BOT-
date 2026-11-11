package store

import (
"encoding/json"
"os"
"sync"
)

type StickerMeta struct {
Pack   string `json:"pack,omitempty"`
Author string `json:"author,omitempty"`
}

type StickerDatabase map[string]StickerMeta

var (
stickerDB   StickerDatabase
stickerMu   sync.Mutex
stickerFile = "database/stickerMeta.json"
)

func init() {
os.MkdirAll("database", 0755)
LoadStickerMeta()
}

func LoadStickerMeta() {
stickerMu.Lock()
defer stickerMu.Unlock()

data, err := os.ReadFile(stickerFile)
if err != nil {
stickerDB = make(StickerDatabase)
return
}

json.Unmarshal(data, &stickerDB)
if stickerDB == nil {
stickerDB = make(StickerDatabase)
}
}

// saveLocked guarda SIN volver a bloquear el mutex (evita deadlock)
func saveLocked() {
data, _ := json.MarshalIndent(stickerDB, "", "  ")
os.WriteFile(stickerFile, data, 0644)
}

func SaveStickerMeta() {
stickerMu.Lock()
defer stickerMu.Unlock()
saveLocked()
}

func GetStickerMeta(jid string) (packname, author string) {
stickerMu.Lock()
defer stickerMu.Unlock()

if meta, exists := stickerDB[jid]; exists {
if meta.Pack != "" {
packname = meta.Pack
} else {
packname = "☄️ Cometa Bot ☄️"
}
if meta.Author != "" {
author = meta.Author
} else {
author = "@" + jid
}
} else {
packname = "☄️ Cometa Bot ☄️"
author = "@" + jid
}

return
}

func SetStickerMeta(jid, pack, auth string) {
stickerMu.Lock()
stickerDB[jid] = StickerMeta{
Pack:   pack,
Author: auth,
}
saveLocked()
stickerMu.Unlock()
}

func HasStickerMeta(jid string) bool {
stickerMu.Lock()
defer stickerMu.Unlock()

if meta, exists := stickerDB[jid]; exists {
return meta.Pack != "" || meta.Author != ""
}
return false
}

func DeleteStickerMeta(jid string) {
stickerMu.Lock()
delete(stickerDB, jid)
saveLocked()
stickerMu.Unlock()
}
