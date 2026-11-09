package cache

import (
"go.mau.fi/whatsmeow"
)

var (
MenuImage   []byte
MenuUpload  *whatsmeow.UploadResponse
MenuGIF     []byte
MenuGIFUpload *whatsmeow.UploadResponse
MenuType    string // "image" o "gif"
)
