package store

import (
"encoding/json"
"os"
"sync"
"time"
)

type Pokemon struct {
ID      int    `json:"id"`
Name    string `json:"name"`
Image   string `json:"image"`
IsShiny bool   `json:"is_shiny"`
Date    string `json:"date"`
Type    string `json:"type"`
}

type PokeUser struct {
Balls     int   `json:"balls"`
LastReset int64 `json:"last_reset"`
}

type PokeSpawn struct {
Pokemon  Pokemon
Attempts int
Cost     int
}

var (
PokemonDB    = make(map[string][]Pokemon)
PokeUsers    = make(map[string]*PokeUser)
PokeSpawns   = make(map[string]*PokeSpawn)
pokemonMu    sync.Mutex
pokemonFile  = "database/pokemon.json"
pokeUserFile = "database/pokeusers.json"
)

const MaxBalls = 3

func init() {
os.MkdirAll("database", 0755)
loadData()
}

func loadData() {
pokemonMu.Lock()
defer pokemonMu.Unlock()

if data, err := os.ReadFile(pokemonFile); err == nil {
json.Unmarshal(data, &PokemonDB)
}
if data, err := os.ReadFile(pokeUserFile); err == nil {
json.Unmarshal(data, &PokeUsers)
}
}

func saveData() {
data1, _ := json.MarshalIndent(PokemonDB, "", "  ")
os.WriteFile(pokemonFile, data1, 0644)
data2, _ := json.MarshalIndent(PokeUsers, "", "  ")
os.WriteFile(pokeUserFile, data2, 0644)
}

func CheckBalls(jid string) (int, bool) {
pokemonMu.Lock()
defer pokemonMu.Unlock()

user, ok := PokeUsers[jid]
now := time.Now().Unix()

if !ok || now-user.LastReset >= 86400 {
PokeUsers[jid] = &PokeUser{Balls: MaxBalls, LastReset: now}
saveData()
return MaxBalls, true
}

return user.Balls, user.Balls > 0
}

func UseBall(jid string) int {
pokemonMu.Lock()
defer pokemonMu.Unlock()

user, ok := PokeUsers[jid]
if !ok {
user = &PokeUser{Balls: MaxBalls, LastReset: time.Now().Unix()}
PokeUsers[jid] = user
}

if user.Balls > 0 {
user.Balls--
}
saveData()
return user.Balls
}

func SetSpawn(jid string, p Pokemon, cost int) {
pokemonMu.Lock()
defer pokemonMu.Unlock()
PokeSpawns[jid] = &PokeSpawn{Pokemon: p, Attempts: 3, Cost: cost}
}

func GetSpawn(jid string) *PokeSpawn {
pokemonMu.Lock()
defer pokemonMu.Unlock()
return PokeSpawns[jid]
}

func ConsumeAttempt(jid string) bool {
pokemonMu.Lock()
defer pokemonMu.Unlock()
spawn := PokeSpawns[jid]
if spawn == nil {
return false
}
spawn.Attempts--
if spawn.Attempts <= 0 {
delete(PokeSpawns, jid)
}
return true
}

func ClearSpawn(jid string) {
pokemonMu.Lock()
defer pokemonMu.Unlock()
delete(PokeSpawns, jid)
}

func AddPokemon(jid string, p Pokemon) {
pokemonMu.Lock()
defer pokemonMu.Unlock()
PokemonDB[jid] = append(PokemonDB[jid], p)
saveData()
}

func GetPokemonInventory(jid string) []Pokemon {
pokemonMu.Lock()
defer pokemonMu.Unlock()
res := make([]Pokemon, len(PokemonDB[jid]))
copy(res, PokemonDB[jid])
return res
}
