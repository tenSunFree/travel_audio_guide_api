package handler

import (
	"fmt"
	"net/http"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) UI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, swaggerHTML)
}

func (h *SwaggerHandler) Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, swaggerSpec)
}

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
  "tags": [{ "name": "Media", "description": "影音刊物" }],
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
            "description": "語系代碼\n\n* zh-tw 正體中文\n* zh-cn 簡體中文\n* en 英文\n* ja 日文\n* ko 韓文",
            "required": true,
            "type": "string"
          },
          {
            "name": "page",
            "in": "query",
            "description": "頁碼（每次回應30筆）",
            "required": false,
            "type": "integer",
            "format": "int32",
            "default": 1
          }
        ],
        "responses": {
          "200": { "description": "OK", "schema": { "$ref": "#/definitions/AudioList" } },
          "400": { "description": "Bad Request" },
          "404": { "description": "Not Found" },
          "502": { "description": "Bad Gateway" }
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
    "AudioList": {
      "type": "object",
      "properties": {
        "total": { "type": "integer", "example": 140 },
        "data":  { "type": "array", "items": { "$ref": "#/definitions/AudioItem" } }
      }
    }
  }
}`

const swaggerHTML = `<!DOCTYPE html>
<html lang="zh-TW">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>旅遊資訊 API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>body { margin: 0; } .topbar { display: none; }</style>
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
      defaultModelsExpandDepth: -1,
      tryItOutEnabled: false
    });
  </script>
</body>
</html>`
