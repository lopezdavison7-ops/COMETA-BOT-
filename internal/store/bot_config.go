package store

import (
"encoding/json"
"os"
"sync"
)

type BotConfig struct {
BannerData []byte `json:"banner_data,omitempty"`
BannerMime string `json:"banner_mime,omitempty"`
BannerType string `json:"banner_type,omitempty"`
BotName    string `json:"bot_name,omitempty"`
}

var (
BotCfg     BotConfig
botCfgMu   sync.Mutex
botCfgFile = "database/bot_config.json"
)

func init() {
os.MkdirAll("database", 0755)
LoadBotConfig()
}

func LoadBotConfig() {
botCfgMu.Lock()
defer botCfgMu.Unlock()
data, err := os.ReadFile(botCfgFile)
if err == nil {
json.Unmarshal(data, &BotCfg)
}
}

func SaveBotConfig() {
botCfgMu.Lock()
defer botCfgMu.Unlock()
data, _ := json.MarshalIndent(BotCfg, "", "  ")
os.WriteFile(botCfgFile, data, 0644)
}

func SetBanner(data []byte, mime, btype string) {
botCfgMu.Lock()
BotCfg.BannerData = data
BotCfg.BannerMime = mime
BotCfg.BannerType = btype
botCfgMu.Unlock()
SaveBotConfig()
}

func SetBotName(name string) {
botCfgMu.Lock()
BotCfg.BotName = name
botCfgMu.Unlock()
SaveBotConfig()
}
