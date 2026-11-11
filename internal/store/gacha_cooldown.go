package store

import (
"encoding/json"
"os"
"sync"
"time"
)

var (
rollCooldowns = make(map[string]int64)
rollMu        sync.Mutex
rollFile      = "database/roll_cooldown.json"
)

const RollCooldownSec = 20 * 60

func init() {
os.MkdirAll("database", 0755)
loadRollCooldowns()
}

func loadRollCooldowns() {
rollMu.Lock()
defer rollMu.Unlock()
data, err := os.ReadFile(rollFile)
if err == nil {
json.Unmarshal(data, &rollCooldowns)
}
}

func saveRollCooldowns() {
rollMu.Lock()
defer rollMu.Unlock()
data, _ := json.Marshal(rollCooldowns)
os.WriteFile(rollFile, data, 0644)
}

func CheckRollCooldown(jid string) (int64, bool) {
rollMu.Lock()
defer rollMu.Unlock()
last, ok := rollCooldowns[jid]
if !ok {
return 0, true
}
now := time.Now().Unix()
if now-last >= RollCooldownSec {
return 0, true
}
return RollCooldownSec - (now - last), false
}

func SetRollCooldown(jid string) {
rollMu.Lock()
rollCooldowns[jid] = time.Now().Unix()
rollMu.Unlock()
saveRollCooldowns()
}
