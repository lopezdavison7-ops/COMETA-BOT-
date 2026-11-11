package store

import (
"encoding/json"
"fmt"
"math/rand"
"os"
"sync"
"time"
)

type GachaCharacter struct {
Name      string `json:"name"`
Series    string `json:"series"`
Gender    string `json:"gender"`
Value     int    `json:"value"`
ImagePath string `json:"image_path"`
Owner     string `json:"owner"`
}

var (
GachaDB    = make(map[string]GachaCharacter)
Inventory  = make(map[string][]string) // jid -> lista de filenames
gachaMu    sync.Mutex
gachaFile  = "database/gacha.json"
inventFile = "database/inventory.json"
)

func init() {
os.MkdirAll("database", 0755)
os.MkdirAll("media/gacha/jpg", 0755)
LoadGacha()
LoadInventory()
}

func LoadGacha() {
gachaMu.Lock()
defer gachaMu.Unlock()
data, err := os.ReadFile(gachaFile)
if err == nil {
json.Unmarshal(data, &GachaDB)
}
}

func SaveGacha() {
gachaMu.Lock()
defer gachaMu.Unlock()
data, _ := json.MarshalIndent(GachaDB, "", "  ")
os.WriteFile(gachaFile, data, 0644)
}

func LoadInventory() {
gachaMu.Lock()
defer gachaMu.Unlock()
data, err := os.ReadFile(inventFile)
if err == nil {
json.Unmarshal(data, &Inventory)
}
}

func SaveInventory() {
gachaMu.Lock()
defer gachaMu.Unlock()
data, _ := json.MarshalIndent(Inventory, "", "  ")
os.WriteFile(inventFile, data, 0644)
}

func GachaCharExists(name, series string) bool {
gachaMu.Lock()
defer gachaMu.Unlock()
for _, c := range GachaDB {
if c.Name == name && c.Series == series {
return true
}
}
return false
}

func AddGachaChar(filename, name, series, gender string, value int) {
gachaMu.Lock()
GachaDB[filename] = GachaCharacter{
Name:      name,
Series:    series,
Gender:    gender,
Value:     value,
ImagePath: "media/gacha/jpg/" + filename,
Owner:     "",
}
gachaMu.Unlock()
SaveGacha()
}

// RollWaifu: selecciona un personaje aleatorio sin dueño
func RollWaifu() (string, GachaCharacter, bool) {
gachaMu.Lock()
defer gachaMu.Unlock()
var available []string
for file, c := range GachaDB {
if c.Owner == "" {
available = append(available, file)
}
}
if len(available) == 0 {
return "", GachaCharacter{}, false
}
r := rand.New(rand.NewSource(time.Now().UnixNano()))
chosen := available[r.Intn(len(available))]
return chosen, GachaDB[chosen], true
}

// Claim: asigna el personaje al usuario
func Claim(jid, filename string) (GachaCharacter, error) {
gachaMu.Lock()
defer gachaMu.Unlock()
c, ok := GachaDB[filename]
if !ok {
return GachaCharacter{}, fmt.Errorf("personaje no existe")
}
if c.Owner != "" {
return GachaCharacter{}, fmt.Errorf("ya tiene dueño")
}
c.Owner = jid
GachaDB[filename] = c
Inventory[jid] = append(Inventory[jid], filename)
gachaMu.Unlock()
SaveGacha()
SaveInventory()
gachaMu.Lock()
return c, nil
}

// GetUserInventory devuelve los personajes de un usuario
func GetUserInventory(jid string) []GachaCharacter {
gachaMu.Lock()
defer gachaMu.Unlock()
var chars []GachaCharacter
for _, file := range Inventory[jid] {
if c, ok := GachaDB[file]; ok {
chars = append(chars, c)
}
}
return chars
}

// SellCharacter: vende un personaje del usuario y devuelve su valor
func SellCharacter(jid, query string) (GachaCharacter, error) {
gachaMu.Lock()
defer gachaMu.Unlock()
files := Inventory[jid]
for i, file := range files {
c, ok := GachaDB[file]
if !ok {
continue
}
// Buscar por nombre o índice
if matchQuery(c.Name, query) {
// Remover del inventario
Inventory[jid] = append(files[:i], files[i+1:]...)
// Liberar el personaje (sin dueño)
c.Owner = ""
GachaDB[file] = c
gachaMu.Unlock()
SaveGacha()
SaveInventory()
gachaMu.Lock()
return c, nil
}
}
return GachaCharacter{}, fmt.Errorf("no tienes ese personaje")
}

// TradeCharacter: intercambia personaje entre dos usuarios
func TradeCharacter(sender, receiver, senderQuery, receiverQuery string) (GachaCharacter, GachaCharacter, error) {
gachaMu.Lock()
defer gachaMu.Unlock()

// Buscar personaje del sender
var senderFile string
var senderChar GachaCharacter
for _, file := range Inventory[sender] {
c, ok := GachaDB[file]
if !ok {
continue
}
if matchQuery(c.Name, senderQuery) {
senderFile = file
senderChar = c
break
}
}
if senderFile == "" {
return GachaCharacter{}, GachaCharacter{}, fmt.Errorf("no tienes ese personaje")
}

// Buscar personaje del receiver
var receiverFile string
var receiverChar GachaCharacter
for _, file := range Inventory[receiver] {
c, ok := GachaDB[file]
if !ok {
continue
}
if matchQuery(c.Name, receiverQuery) {
receiverFile = file
receiverChar = c
break
}
}
if receiverFile == "" {
return GachaCharacter{}, GachaCharacter{}, fmt.Errorf("el otro usuario no tiene ese personaje")
}

// Hacer el intercambio
senderChar.Owner = receiver
receiverChar.Owner = sender
GachaDB[senderFile] = senderChar
GachaDB[receiverFile] = receiverChar

// Actualizar inventarios
removeFromSlice := func(s []string, v string) []string {
for i, x := range s {
if x == v {
return append(s[:i], s[i+1:]...)
}
}
return s
}
Inventory[sender] = removeFromSlice(Inventory[sender], senderFile)
Inventory[sender] = append(Inventory[sender], receiverFile)
Inventory[receiver] = removeFromSlice(Inventory[receiver], receiverFile)
Inventory[receiver] = append(Inventory[receiver], senderFile)

gachaMu.Unlock()
SaveGacha()
SaveInventory()
gachaMu.Lock()

return senderChar, receiverChar, nil
}

func matchQuery(name, query string) bool {
if query == "" {
return false
}
// Búsqueda simple: contiene o igual (case insensitive)
nameLower := ""
for _, r := range name {
if r >= 'A' && r <= 'Z' {
nameLower += string(r + 32)
} else {
nameLower += string(r)
}
}
queryLower := ""
for _, r := range query {
if r >= 'A' && r <= 'Z' {
queryLower += string(r + 32)
} else {
queryLower += string(r)
}
}
return nameLower == queryLower || contains(nameLower, queryLower)
}

func contains(s, sub string) bool {
if len(sub) > len(s) {
return false
}
for i := 0; i <= len(s)-len(sub); i++ {
if s[i:i+len(sub)] == sub {
return true
}
}
return false
}
