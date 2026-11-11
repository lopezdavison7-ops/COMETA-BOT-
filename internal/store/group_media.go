package store

import (
"encoding/json"
"os"
"path/filepath"
"sync"

"go.mau.fi/whatsmeow/proto/waE2E"
"google.golang.org/protobuf/proto"
)

type SavedMedia struct {
URL           string `json:"url"`
DirectPath    string `json:"direct_path"`
MediaKey      []byte `json:"media_key"`
FileEncSHA256 []byte `json:"file_enc_sha256"`
FileSHA256    []byte `json:"file_sha256"`
FileLength    uint64 `json:"file_length"`
Mimetype      string `json:"mimetype"`
}

var (
mediaMu   sync.Mutex
mediaPath = "database/group_media"
)

func init() {
os.MkdirAll(mediaPath, 0755)
}

func getMediaPath(jid, action string) string {
safeJID := filepath.Base(jid)
return filepath.Join(mediaPath, safeJID+"_"+action+".json")
}

func SaveMediaUpload(jid, action string, media SavedMedia) error {
mediaMu.Lock()
defer mediaMu.Unlock()
data, err := json.Marshal(media)
if err != nil {
return err
}
return os.WriteFile(getMediaPath(jid, action), data, 0644)
}

func GetMediaUpload(jid, action string) *waE2E.ImageMessage {
mediaMu.Lock()
defer mediaMu.Unlock()

data, err := os.ReadFile(getMediaPath(jid, action))
if err != nil {
return nil
}

var saved SavedMedia
if err := json.Unmarshal(data, &saved); err != nil {
return nil
}

return &waE2E.ImageMessage{
URL:           proto.String(saved.URL),
DirectPath:    proto.String(saved.DirectPath),
MediaKey:      saved.MediaKey,
FileEncSHA256: saved.FileEncSHA256,
FileSHA256:    saved.FileSHA256,
FileLength:    proto.Uint64(saved.FileLength),
Mimetype:      proto.String(saved.Mimetype),
}
}

func DeleteMediaUpload(jid, action string) {
mediaMu.Lock()
defer mediaMu.Unlock()
os.Remove(getMediaPath(jid, action))
}
