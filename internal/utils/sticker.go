package utils

import (
"crypto/rand"
"encoding/hex"
"fmt"
"os"
"os/exec"
"path/filepath"
)

const BOT_PACK = "☄️ Cometa Bot ☄️"
const BOT_AUTHOR = "© DevLuis"

func init() {
os.MkdirAll("tmp", 0755)
}

func randomID() string {
b := make([]byte, 8)
rand.Read(b)
return hex.EncodeToString(b)
}

func ConvertImageToWebp(buffer []byte) ([]byte, error) {
id := randomID()
inputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s.jpg", id))
outputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s.webp", id))

defer os.Remove(inputPath)
defer os.Remove(outputPath)

if err := os.WriteFile(inputPath, buffer, 0644); err != nil {
return nil, err
}

cmd := exec.Command("ffmpeg",
"-i", inputPath,
"-vf", "scale=512:512:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=0x00000000",
"-c:v", "libwebp",
"-lossless", "0",
"-q:v", "80",
outputPath,
"-y",
)

if err := cmd.Run(); err != nil {
return nil, fmt.Errorf("error convirtiendo imagen: %v", err)
}

return os.ReadFile(outputPath)
}

func ConvertVideoToWebp(buffer []byte) ([]byte, error) {
id := randomID()
inputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s.mp4", id))
outputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s.webp", id))

defer os.Remove(inputPath)
defer os.Remove(outputPath)

if err := os.WriteFile(inputPath, buffer, 0644); err != nil {
return nil, err
}

durationCmd := exec.Command("ffprobe",
"-v", "error",
"-show_entries", "format=duration",
"-of", "default=noprint_wrappers=1:nokey=1",
inputPath,
)

durationOut, err := durationCmd.Output()
if err != nil {
return nil, fmt.Errorf("error leyendo duración: %v", err)
}

var duration float64
fmt.Sscanf(string(durationOut), "%f", &duration)

if duration > 15 {
return nil, fmt.Errorf("el video dura %.1fs, máximo 15 segundos", duration)
}

cmd := exec.Command("ffmpeg",
"-i", inputPath,
"-t", "15",
"-vf", "fps=15,scale=512:512:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=0x00000000",
"-c:v", "libwebp",
"-lossless", "0",
"-q:v", "70",
"-loop", "0",
"-an",
"-vsync", "0",
outputPath,
"-y",
)

if err := cmd.Run(); err != nil {
return nil, fmt.Errorf("error convirtiendo video: %v", err)
}

return os.ReadFile(outputPath)
}

func CleanSticker(buffer []byte) ([]byte, error) {
id := randomID()
inputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s.webp", id))
outputPath := filepath.Join("tmp", fmt.Sprintf("stk_%s_out.webp", id))

defer os.Remove(inputPath)
defer os.Remove(outputPath)

if err := os.WriteFile(inputPath, buffer, 0644); err != nil {
return nil, err
}

cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "copy", outputPath, "-y")
if err := cmd.Run(); err != nil {
cmd = exec.Command("ffmpeg",
"-i", inputPath,
"-vf", "scale=512:512:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=0x00000000",
"-c:v", "libwebp",
"-lossless", "0",
"-q:v", "70",
"-loop", "0",
"-an",
"-vsync", "0",
outputPath,
"-y",
)
if err := cmd.Run(); err != nil {
return buffer, nil
}
}

return os.ReadFile(outputPath)
}
