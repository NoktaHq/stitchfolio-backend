// patch_swagger_security.go injects root-level BearerAuth security into swag-generated
// docs so Swagger UI sends the Token header. Run after: swag init
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	docsDir := "docs"
	if len(os.Args) > 1 {
		docsDir = os.Args[1]
	}

	// 1. docs/docs.go
	docsGoPath := filepath.Join(docsDir, "docs.go")
	if err := patchDocsGo(docsGoPath); err != nil {
		fmt.Fprintf(os.Stderr, "docs.go: %v\n", err)
		os.Exit(1)
	}

	// 2. docs/swagger.json
	swaggerJSONPath := filepath.Join(docsDir, "swagger.json")
	if err := patchSwaggerJSON(swaggerJSONPath); err != nil {
		fmt.Fprintf(os.Stderr, "swagger.json: %v\n", err)
		os.Exit(1)
	}

	// 3. docs/swagger.yaml
	swaggerYAMLPath := filepath.Join(docsDir, "swagger.yaml")
	if err := patchSwaggerYAML(swaggerYAMLPath); err != nil {
		fmt.Fprintf(os.Stderr, "swagger.yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Patched swagger docs with root-level BearerAuth security.")
}

func patchDocsGo(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(body)

	// Already has security block?
	if strings.Contains(content, `"BearerAuth": []`) {
		return nil
	}

	// Insert security between basePath and paths in the docTemplate
	old := `    "basePath": "{{.BasePath}}",
    "paths": {`
	new := `    "basePath": "{{.BasePath}}",
    "security": [
        {
            "BearerAuth": []
        }
    ],
    "paths": {`
	if !strings.Contains(content, old) {
		return fmt.Errorf("expected pattern not found in docs.go (swag may have changed output)")
	}
	content = strings.Replace(content, old, new, 1)
	return os.WriteFile(path, []byte(content), 0644)
}

func patchSwaggerJSON(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(body)

	if strings.Contains(content, `"BearerAuth": []`) {
		return nil
	}

	// Match: "basePath": "<any>",\n    "paths":
	re := regexp.MustCompile(`("basePath":\s*"[^"]*",)\s*("paths":\s*\{)`)
	repl := "$1\n    \"security\": [\n        {\n            \"BearerAuth\": []\n        }\n    ],\n    $2"
	newContent := re.ReplaceAllString(content, repl)
	if newContent == content {
		return fmt.Errorf("expected pattern not found in swagger.json")
	}
	return os.WriteFile(path, []byte(newContent), 0644)
}

func patchSwaggerYAML(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(body)

	if strings.Contains(content, "BearerAuth: []") {
		return nil
	}

	// basePath: ... followed by definitions: (allow optional blank line)
	re := regexp.MustCompile(`(basePath:\s*/api/sf/v1)\s*\n(definitions:)`)
	repl := "${1}\nsecurity:\n  - BearerAuth: []\n${2}"
	newContent := re.ReplaceAllString(content, repl)
	if newContent == content {
		return fmt.Errorf("expected pattern not found in swagger.yaml")
	}
	return os.WriteFile(path, []byte(newContent), 0644)
}
