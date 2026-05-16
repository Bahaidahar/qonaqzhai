package handler

import (
	"embed"
	"net/http"
)

//go:embed docs/*
var docsFS embed.FS

const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Qonaqzhai API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/api/docs/openapi.yaml",
        dom_id: "#swagger",
      });
    };
  </script>
</body>
</html>`

// SwaggerUI serves an HTML page that loads the OpenAPI spec.
func SwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerUI))
}

// OpenAPI serves the embedded openapi.yaml.
func OpenAPI(w http.ResponseWriter, _ *http.Request) {
	body, err := docsFS.ReadFile("docs/openapi.yaml")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(body)
}
