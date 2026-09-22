// dav-drive-record: TasiaSnap Recorder (beta) for Windows.
// Records a short screen clip via ffmpeg (gdigrab) and uploads it to
// Dav Drive through the ShareX-compatible endpoint.
//
// Usage: dav-drive-record.exe [-d 15]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func loadConf(dir string) map[string]string {
	m := map[string]string{}
	candidates := []string{filepath.Join(dir, "dav-drive.conf")}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".config", "TasiaSnap", "dav-drive.conf"))
	}
	var f *os.File
	var err error
	for _, c := range candidates {
		f, err = os.Open(c)
		if err == nil {
			break
		}
	}
	if f == nil {
		return m
	}
	defer f.Close()
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.IndexByte(line, '='); i > 0 {
			m[strings.TrimSpace(line[:i])] = trimQuotes(strings.TrimSpace(line[i+1:]))
		}
	}
	return m
}

func trimQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"') {
		return s[1 : len(s)-1]
	}
	return s
}

func uploadFile(host, token, folder, path, filename string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf strings.Builder
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("token", token)
	_ = mw.WriteField("folder", folder)
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, file); err != nil {
		return "", err
	}
	if err := mw.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", host+"/upload/sharex", strings.NewReader(buf.String()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	key := `"url"`
	idx := strings.Index(string(body), key)
	if idx < 0 {
		return "", fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
	}
	rest := string(body)[idx+len(key):]
	open := strings.IndexByte(rest, '"')
	if open < 0 {
		return "", fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
	}
	val := rest[open+1:]
	if close := strings.IndexByte(val, '"'); close >= 0 {
		val = val[:close]
	}
	if !strings.HasPrefix(val, "http") {
		return "", fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
	}
	return val, nil
}

func main() {
	dur := flag.Int("d", 15, "recording length in seconds (keep it short)")
	flag.Parse()

	fmt.Println("== TasiaSnap Recorder (beta) ==")
	fmt.Printf("Recording %d seconds... keep it short, and please expect rough edges.\n\n", *dur)

	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ffmpeg not found.")
		fmt.Fprintln(os.Stderr, "Install it from https://www.gyan.dev/ffmpeg/builds/ and add ffmpeg.exe to PATH.")
		os.Exit(1)
	}

	out := filepath.Join(os.TempDir(), fmt.Sprintf("tasiasnap-record-%d.mp4", time.Now().Unix()))

	cmd := exec.Command(ffmpeg, "-y", "-loglevel", "error",
		"-f", "gdigrab", "-framerate", "20", "-i", "desktop",
		"-t", fmt.Sprintf("%d", *dur),
		"-c:v", "libx264", "-pix_fmt", "yuv420p", out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: recording failed:", err)
		os.Exit(1)
	}

	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cfg := loadConf(exeDir)
	host := cfg["DAV_DRIVE_HOST"]
	token := cfg["DAV_DRIVE_TOKEN"]
	folder := cfg["DAV_DRIVE_FOLDER"]
	if folder == "" {
		folder = "/ShareX/"
	}
	if host == "" || token == "" {
		fmt.Fprintln(os.Stderr, "dav-drive.conf must define DAV_DRIVE_HOST and DAV_DRIVE_TOKEN")
		os.Exit(1)
	}

	fmt.Println("Uploading to Dav Drive...")
	url, err := uploadFile(host, token, folder, out, filepath.Base(out))
	_ = os.Remove(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Clip link:", url)
	fmt.Println("(Recorder is beta - short recordings advised.)")
}