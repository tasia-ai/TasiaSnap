// dav-drive-uploader: Dav Drive uploader for Ksnip "Script Uploader" (Windows).
// Ksnip passes the screenshot path as argv[1]; this helper posts it to Dav
// Drive via the ShareX-compatible endpoint and prints the public URL on stdout.
package main

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dav-drive-uploader.exe <image-path>")
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

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	var buf strings.Builder
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("token", token)
	_ = mw.WriteField("folder", folder)
	fw, err := mw.CreateFormFile("file", filepath.Base(os.Args[1]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := io.Copy(fw, file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := mw.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", host+"/upload/sharex", strings.NewReader(buf.String()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	key := `"url"`
	idx := strings.Index(string(body), key)
	if idx < 0 {
		fmt.Fprintf(os.Stderr, "upload failed (%d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}
	rest := string(body)[idx+len(key):]
	open := strings.IndexByte(rest, '"')
	if open < 0 {
		fmt.Fprintf(os.Stderr, "upload failed (%d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}
	val := rest[open+1:]
	if close := strings.IndexByte(val, '"'); close >= 0 {
		val = val[:close]
	}
	if !strings.HasPrefix(val, "http") {
		fmt.Fprintf(os.Stderr, "upload failed (%d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}
	fmt.Println(val)
}