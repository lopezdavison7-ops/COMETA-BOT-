package store

import "sort"

type CollectorStats struct {
JID        string
Count      int
TotalValue int
Best       GachaCharacter
}

func GetTopCollectors(limit int) []CollectorStats {
gachaMu.Lock()
defer gachaMu.Unlock()

stats := make(map[string]*CollectorStats)
for _, c := range GachaDB {
if c.Owner == "" {
continue
}
s, ok := stats[c.Owner]
if !ok {
s = &CollectorStats{JID: c.Owner}
stats[c.Owner] = s
}
s.Count++
s.TotalValue += c.Value
if s.Best.Name == "" || c.Value > s.Best.Value {
s.Best = c
}
}

var list []CollectorStats
for _, s := range stats {
list = append(list, *s)
}

sort.Slice(list, func(i, j int) bool {
if list[i].Count != list[j].Count {
return list[i].Count > list[j].Count
}
return list[i].TotalValue > list[j].TotalValue
})

if len(list) > limit {
list = list[:limit]
}
return list
}

func GetTopWaifus(limit int) []GachaCharacter {
gachaMu.Lock()
defer gachaMu.Unlock()
var owned []GachaCharacter
for _, c := range GachaDB {
if c.Owner != "" {
owned = append(owned, c)
}
}
sort.Slice(owned, func(i, j int) bool {
return owned[i].Value > owned[j].Value
})
if len(owned) > limit {
owned = owned[:limit]
}
return owned
}

func GetUserTotalValue(jid string) int {
gachaMu.Lock()
defer gachaMu.Unlock()
total := 0
for _, file := range Inventory[jid] {
if c, ok := GachaDB[file]; ok {
total += c.Value
}
}
return total
}
