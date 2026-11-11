package gacha

import (
"encoding/json"
"fmt"
"io"
"math/rand"
"net/http"
"net/url"
"os"
"strings"
"time"

"cometa-bot/internal/store"
)

var httpClient = &http.Client{Timeout: 20 * time.Second}

var hosts = []string{"https://konachan.net", "https://konachan.com"}

const userAgent = "konachan-scraper/1.0"

var bannedTags = map[string]bool{
"loli": true, "shota": true, "child": true, "toddler": true, "infant": true,
"nude": true, "naked": true, "topless": true, "uncensored": true, "censored": true,
"pussy": true, "penis": true, "cum": true, "sex": true, "fellatio": true,
}

var skipTags = map[string]bool{
"highres": true, "absurdres": true, "jpeg_artifacts": true, "scan": true, "dakimakura": true,
"1girl": true, "2girls": true, "3girls": true, "4girls": true, "multiple_girls": true, "solo": true,
"1boy": true, "2boys": true, "multiple_boys": true, "swimsuits": true, "thighhighs": true,
"bikini": true, "wet": true, "pantsu": true, "nipples": true, "dress": true, "see_through": true,
"animal_ears": true, "ass": true, "skirt_lift": true, "open_shirt": true, "bra": true, "tail": true,
"breasts": true, "cleavage": true, "panties": true, "navel": true, "blush": true, "long_hair": true,
"short_hair": true, "blonde_hair": true, "twintails": true, "brown_hair": true, "black_hair": true,
"white_hair": true, "red_hair": true, "blue_hair": true, "green_hair": true, "no_bra": true,
"megane": true, "horns": true, "stockings": true, "pantyhose": true, "weapon": true, "cosplay": true,
"bunny_ears": true, "feet": true, "lingerie": true, "bunny_girl": true, "leotard": true, "sword": true,
"armor": true, "torn_clothes": true, "seifuku": true, "wings": true, "shirt_lift": true,
"wedding_dress": true, "gym_uniform": true, "maid": true, "towel": true, "naked_apron": true,
"yukata": true, "uniform": true, "pajama": true, "underboob": true, "shimapan": true,
"vector_trace": true, "wallpaper": true, "transparent_png": true, "monochrome": true,
"crossover": true, "tagme": true, "fixme": true, "crease": true, "onsen": true, "yuri": true,
}

var valoresPosibles = []int{1200, 2500, 3800, 4660, 5200, 6800, 8500, 10000, 12400, 15000, 18000, 25000}

func doGet(u string) (*http.Response, error) {
req, err := http.NewRequest("GET", u, nil)
if err != nil {
return nil, err
}
req.Header.Set("User-Agent", userAgent)
req.Header.Set("Accept", "application/json, text/plain, */*")
return httpClient.Do(req)
}

func fetchJSON(path string, target interface{}) error {
var lastErr error
for _, host := range hosts {
resp, err := doGet(host + path)
if err != nil {
lastErr = err
fmt.Printf("⚠️ Gacha: %s error de red: %v\n", host, err)
continue
}
if resp.StatusCode != 200 {
lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
fmt.Printf("⚠️ Gacha: %s devolvió HTTP %d\n", host, resp.StatusCode)
resp.Body.Close()
continue
}
err = json.NewDecoder(resp.Body).Decode(target)
resp.Body.Close()
if err != nil {
lastErr = err
fmt.Printf("⚠️ Gacha: %s JSON inválido: %v\n", host, err)
continue
}
return nil
}
return lastErr
}

func RandomValue() int {
return valoresPosibles[rand.Intn(len(valoresPosibles))]
}

func TagToName(tag string) string {
tag = strings.ReplaceAll(tag, "_", " ")
words := strings.Fields(tag)
for i, w := range words {
if len(w) > 0 {
words[i] = strings.ToUpper(w[:1]) + w[1:]
}
}
return strings.Join(words, " ")
}

func FetchPosts(seriesTag string, extraTags []string, pages int) []map[string]interface{} {
baseTags := append([]string{seriesTag}, extraTags...)
tagStr := strings.Join(baseTags, " ")
var allPosts []map[string]interface{}

for page := 1; page <= pages; page++ {
path := fmt.Sprintf("/post.json?tags=%s&limit=100&page=%d", url.QueryEscape(tagStr), page)
var posts []map[string]interface{}
if err := fetchJSON(path, &posts); err != nil {
break
}
if len(posts) == 0 {
break
}

for _, p := range posts {
tagsStr, _ := p["tags"].(string)
rating, _ := p["rating"].(string)
if rating == "e" {
continue
}
isBanned := false
for t := range bannedTags {
if strings.Contains(tagsStr, t) {
isBanned = true
break
}
}
if !isBanned {
allPosts = append(allPosts, p)
}
}
if len(posts) < 100 {
break
}
time.Sleep(800 * time.Millisecond)
}
return allPosts
}

func GetCharacterTags(posts []map[string]interface{}, seriesTag string) []string {
freq := make(map[string]int)
for _, p := range posts {
tagsStr, _ := p["tags"].(string)
for _, t := range strings.Fields(tagsStr) {
if skipTags[t] || t == seriesTag {
continue
}
freq[t]++
}
}

var candidates []string
for t, count := range freq {
if count >= 2 {
candidates = append(candidates, t)
}
}

var characters []string
for _, tag := range candidates {
path := fmt.Sprintf("/tag.json?name=%s", url.QueryEscape(tag))
var tagData []map[string]interface{}
if err := fetchJSON(path, &tagData); err != nil {
continue
}

isChar := false
for _, td := range tagData {
name, _ := td["name"].(string)
ttype, ok := td["type"].(float64)
if name == tag && ok && int(ttype) == 4 {
isChar = true
break
}
}

if isChar {
characters = append(characters, tag)
}
time.Sleep(500 * time.Millisecond)
}
return characters
}

func GetGender(charTag string, posts []map[string]interface{}) string {
if strings.Contains(charTag, "_(male)") {
return "Masculino"
}
if strings.Contains(charTag, "_(female)") {
return "Femenino"
}

femaleTags := map[string]bool{"1girl": true, "2girls": true, "multiple_girls": true, "female": true}
maleTags := map[string]bool{"1boy": true, "2boys": true, "multiple_boys": true, "male": true}

fScore, mScore := 0, 0
for _, p := range posts {
tagsStr, _ := p["tags"].(string)
if !strings.Contains(tagsStr, charTag) {
continue
}
tags := strings.Fields(tagsStr)
weight := 1
if len(tags) > 0 && tags[0] == "solo" {
weight = 10
}
for _, t := range tags {
if femaleTags[t] {
fScore += weight
}
if maleTags[t] {
mScore += weight
}
}
}

if mScore > fScore {
return "Masculino"
}
return "Femenino"
}

func DownloadAndSaveImage(post map[string]interface{}, charName, seriesName string) (string, error) {
sampleURL, _ := post["sample_url"].(string)
fileURL, _ := post["file_url"].(string)
imgURL := sampleURL
if imgURL == "" {
imgURL = fileURL
}
if imgURL == "" {
return "", fmt.Errorf("sin URL de imagen")
}

resp, err := doGet(imgURL)
if err != nil || resp.StatusCode != 200 {
return "", fmt.Errorf("error descargando imagen")
}
defer resp.Body.Close()

data, _ := io.ReadAll(resp.Body)
if len(data) == 0 {
return "", fmt.Errorf("imagen vacía")
}

filename := fmt.Sprintf("gacha_%d_%d.jpg", time.Now().Unix(), rand.Intn(10000))
filepath := "media/gacha/jpg/" + filename

os.WriteFile(filepath, data, 0644)

store.AddGachaChar(filename, charName, seriesName, GetGender(charName, []map[string]interface{}{post}), RandomValue())

return filename, nil
}

func FetchRandomSeriesTags(cantidad int) []string {
maxPage := 15
minCount := 30
pool := make(map[string]bool)

pages := make(map[int]bool)
for len(pages) < 3 {
pages[rand.Intn(maxPage)+1] = true
}

for page := range pages {
path := fmt.Sprintf("/tag.json?type=3&order=count&limit=100&page=%d", page)
var tags []map[string]interface{}
if err := fetchJSON(path, &tags); err != nil {
continue
}
for _, t := range tags {
count, _ := t["count"].(float64)
name, _ := t["name"].(string)
if int(count) >= minCount && name != "" {
pool[name] = true
}
}
time.Sleep(400 * time.Millisecond)
}

if len(pool) == 0 {
var tags []map[string]interface{}
if err := fetchJSON("/tag.json?type=3&order=count&limit=100&page=1", &tags); err == nil {
for _, t := range tags {
count, _ := t["count"].(float64)
name, _ := t["name"].(string)
if int(count) >= minCount && name != "" {
pool[name] = true
}
}
}
}

result := make([]string, 0, len(pool))
for tag := range pool {
result = append(result, tag)
}
rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
if len(result) > cantidad {
result = result[:cantidad]
}
return result
}
