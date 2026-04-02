package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ─────────────────────────────────────────────
// Models
// ─────────────────────────────────────────────

type AudioItem struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Summary  *string `json:"summary"`
	URL      string  `json:"url"`
	FileExt  *string `json:"file_ext"`
	Modified string  `json:"modified"`
}

type AudioResponse struct {
	Total int         `json:"total"`
	Data  []AudioItem `json:"data"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ─────────────────────────────────────────────
// Supported languages
// ─────────────────────────────────────────────

var validLangs = map[string]bool{
	"zh-tw": true,
	"zh-cn": true,
	"en":    true,
	"ja":    true,
	"ko":    true,
}

// Upstream real API
const upstreamBase = "https://www.travel.taipei/open-api"

// ─────────────────────────────────────────────
// Middleware
// ─────────────────────────────────────────────

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s | %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ─────────────────────────────────────────────
// Handler: Proxy requests to the real travel.taipei API
// ─────────────────────────────────────────────

func audioHandler(w http.ResponseWriter, r *http.Request) {
	// Expected path format: /open-api/{lang}/Media/Audio
	// parts[0]="open-api"  parts[1]=lang  parts[2]="Media"  parts[3]="Audio"
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeJSON(w, http.StatusNotFound, ErrorResponse{404, "Not Found"})
		return
	}

	lang := parts[1]
	if !validLangs[lang] {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{400, fmt.Sprintf("不支援的語系: %s", lang)})
		return
	}

	// Parse page number, default to 1
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Proxy the request to the upstream travel.taipei API
	upstreamURL := fmt.Sprintf("%s/%s/Media/Audio?page=%d", upstreamBase, lang, page)
	log.Printf("→ proxy: %s", upstreamURL)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", upstreamURL, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{500, "建立請求失敗"})
		return
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{500, "上游 API 請求失敗: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{500, "讀取回應失敗"})
		return
	}

	// Forward the upstream response as-is
	// This preserves real values such as null fields and total: 140
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// ─────────────────────────────────────────────
// Swagger Handlers
// ─────────────────────────────────────────────

func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, swaggerHTML)
}

func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, swaggerSpec)
}

// ─────────────────────────────────────────────
// Router
// ─────────────────────────────────────────────

func newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		// /open-api/{lang}/Media/Audio
		if len(parts) == 4 &&
			parts[0] == "open-api" &&
			strings.EqualFold(parts[2], "Media") &&
			strings.EqualFold(parts[3], "Audio") {
			audioHandler(w, r)
			return
		}

		// Swagger UI
		if path == "" ||
			path == "open-api/swagger/ui" ||
			path == "open-api/swagger/ui/index" {
			swaggerUIHandler(w, r)
			return
		}

		// Swagger JSON
		if path == "open-api/swagger/docs" ||
			path == "open-api/swagger/docs/V1" {
			swaggerJSONHandler(w, r)
			return
		}

		writeJSON(w, http.StatusNotFound, ErrorResponse{404, "Not Found"})
	})

	return loggingMiddleware(corsMiddleware(mux))
}

// ─────────────────────────────────────────────
// Main
// ─────────────────────────────────────────────

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("🚀 台北語音導覽 API 啟動中...")
	log.Printf("📡 API:          http://localhost%s/open-api/{lang}/Media/Audio", addr)
	log.Printf("📖 Swagger UI:   http://localhost%s/open-api/swagger/ui/index", addr)
	log.Printf("📄 Swagger JSON: http://localhost%s/open-api/swagger/docs/V1", addr)

	if err := http.ListenAndServe(addr, newRouter()); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}

// ─────────────────────────────────────────────
// Swagger OpenAPI 2.0 Specification
// ─────────────────────────────────────────────
// Remove "default" and "enum" so Swagger UI shows a normal text input
const swaggerSpec = `{
  "swagger": "2.0",
  "info": {
    "title": "旅遊資訊API",
    "description": "本開放資料平台透過swagger套件提供。",
    "version": "V1"
  },
  "host": "localhost:8080",
  "basePath": "/open-api",
  "schemes": ["http"],
  "tags": [
    {
      "name": "Media",
      "description": "影音刊物"
    }
  ],
  "paths": {
    "/{lang}/Media/Audio": {
      "get": {
        "tags": ["Media"],
        "summary": "語音導覽",
        "operationId": "Media_Audio",
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "lang",
            "in": "path",
            "description": "語系代碼\n\n* zh-tw -正體中文\n* zh-cn -簡體中文\n* en -英文\n* ja -日文\n* ko -韓文",
            "required": true,
            "type": "string",
            // "default": "zh-tw",
            // "enum": ["zh-tw", "zh-cn", "en", "ja", "ko"]
          },
          {
            "name": "page",
            "in": "query",
            "description": "頁碼。(每次回應30筆資料)",
            "required": false,
            "type": "integer",
            "format": "int32",
            "default": 1
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/AudioResponse" }
          },
          "204": { "description": "No Content" },
          "403": { "description": "Forbidden" },
          "404": { "description": "Not Found" },
          "500": { "description": "System Busy" }
        }
      }
    }
  },
  "definitions": {
    "AudioItem": {
      "type": "object",
      "properties": {
        "id":       { "type": "integer", "example": 28 },
        "title":    { "type": "string",  "example": "北投圖書館" },
        "summary":  { "type": "string",  "x-nullable": true },
        "url":      { "type": "string",  "example": "https://www.travel.taipei/audio/28" },
        "file_ext": { "type": "string",  "x-nullable": true },
        "modified": { "type": "string",  "example": "2025-12-10 15:55:41 +08:00" }
      }
    },
    "AudioResponse": {
      "type": "object",
      "properties": {
        "total": { "type": "integer", "example": 140 },
        "data": {
          "type": "array",
          "items": { "$ref": "#/definitions/AudioItem" }
        }
      }
    }
  }
}`

// ─────────────────────────────────────────────
// Swagger UI HTML (CDN)
// ─────────────────────────────────────────────

const swaggerHTML = `<!DOCTYPE html>
<html lang="zh-TW">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>旅遊資訊 API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body { margin: 0; background: #fafafa; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/open-api/swagger/docs/V1",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout",
      deepLinking: true,
      displayOperationId: false,
	  // Hide Models section
      defaultModelsExpandDepth: -1,
	  // Expand model fields by one level
      defaultModelExpandDepth: -1,
	  // Do not enable "Try it out" automatically
      tryItOutEnabled: false
    });
  </script>
</body>
</html>`
