package store

import (
"encoding/json"
"os"
"sync"
"time"
)

type UserProfile struct {
Age        int    `json:"age"`
Gender     string `json:"gender"`
Desc       string `json:"desc"`
FavSong    string `json:"fav_song"`
Crush      string `json:"crush"`
FirstSeen  int64  `json:"first_seen"`
}

var (
profiles   = make(map[string]*UserProfile)
profileMu  sync.Mutex
profileFile = "database/profiles.json"
)

func init() {
os.MkdirAll("database", 0755)
loadProfiles()
}

func loadProfiles() {
profileMu.Lock()
defer profileMu.Unlock()
data, err := os.ReadFile(profileFile)
if err == nil {
json.Unmarshal(data, &profiles)
}
}

func saveProfiles() {
data, _ := json.MarshalIndent(profiles, "", "  ")
os.WriteFile(profileFile, data, 0644)
}

func GetProfile(jid string) *UserProfile {
profileMu.Lock()
defer profileMu.Unlock()
p, ok := profiles[jid]
if !ok {
p = &UserProfile{FirstSeen: time.Now().Unix()}
profiles[jid] = p
saveProfiles()
}
return p
}

func SetProfileField(jid, field, value string) {
profileMu.Lock()
defer profileMu.Unlock()
p, ok := profiles[jid]
if !ok {
p = &UserProfile{FirstSeen: time.Now().Unix()}
profiles[jid] = p
}
switch field {
case "age":
var age int
json.Unmarshal([]byte(value), &age)
p.Age = age
case "gender":
p.Gender = value
case "desc":
p.Desc = value
case "fav_song":
p.FavSong = value
case "crush":
p.Crush = value
}
saveProfiles()
}
