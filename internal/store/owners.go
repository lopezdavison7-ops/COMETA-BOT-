package store

import (
"strings"
"sync"
)

var (
ownersList = []string{
"50588910385",
"50578391933",
"50577455383",
}

lidToPhone = make(map[string]string)
ownersMu   sync.RWMutex
)

func IsOwner(jid string) bool {
ownersMu.RLock()
defer ownersMu.RUnlock()

num := strings.Split(jid, "@")[0]

for _, owner := range ownersList {
if num == owner {
return true
}
}

if phone, ok := lidToPhone[jid]; ok {
for _, owner := range ownersList {
if phone == owner {
return true
}
}
}

return false
}

func RegisterLID(lid, phone string) {
if lid == "" || phone == "" {
return
}
ownersMu.Lock()
lidToPhone[lid] = phone
ownersMu.Unlock()
}

func GetPhoneFromLID(lid string) (string, bool) {
ownersMu.RLock()
defer ownersMu.RUnlock()
phone, ok := lidToPhone[lid]
return phone, ok
}
