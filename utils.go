package main

import (
	"io"
	"log"
	"net"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseRequestLine(requestLine string) (method, path string) {
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) >= 2 {
		method = parts[0]
		path = parts[1]
	}
	return
}

func getContentType(filePath string) string {
	// ファイルの拡張子に基づいてMIMEタイプを決定
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream" // 不明なファイルタイプ
	}
}

func serveFile(conn net.Conn, path string, clientAddr string) {
	/* For directory traversal */
	cleanPath := filepath.Clean("/" + path)
	if !strings.HasPrefix(cleanPath, "/") {
		io.WriteString(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
		log.Printf("Addr: %s, Bad Request: %s\n", clientAddr, path)
		return
	}

	filePath := filepath.Join(baseDir, cleanPath)
	if _, err := filepath.Rel(baseDir, filePath); err != nil || strings.Contains(filePath, "..") {
		io.WriteString(conn, "HTTP/1.1 403 Forbidden\r\n\r\n")
		log.Printf("Addr: %s, Forbidden Path: %s\n", clientAddr, path)
		return
	}
	/* For directory traversal */

	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}

	file, err := os.Open(filePath)
	if err != nil {
		io.WriteString(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
		log.Printf("File Not Found: %s\n", filePath)
		return
	}
	defer file.Close()

	contentType := getContentType(filePath)

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	io.WriteString(conn, fmt.Sprintf("Content-Type: %s\r\n", contentType))
	io.WriteString(conn, "\r\n")
	io.Copy(conn, file)
	log.Printf("Addr: %s, Served %s\n", clientAddr, filePath)
	
}