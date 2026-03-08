package bms_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/app"
	"github.com/pafthang/bms/internal/config"
)

func TestOpenAPIContractQuality(t *testing.T) {
	cfg := config.Config{
		HTTPAddr:        ":0",
		ShutdownTimeout: 3 * time.Second,
		AllowedOrigins:  []string{"http://localhost:3000"},
		DBDSN:           "file::memory:?_pragma=foreign_keys(1)",
		JWTSecret:       "test-secret",
		AccessTTL:       time.Hour,
		RefreshTTL:      24 * time.Hour,
	}
	e := app.BuildEngine(cfg, nil, app.BuildOptions{})
	spec := e.OpenAPISpec()
	if spec == nil {
		t.Fatal("openapi spec is nil")
	}

	issues := arc.ValidateOpenAPIQuality(spec, arc.OpenAPIQualityGates{
		RequireRootTags:         true,
		RequireServers:          true,
		RequireExamples:         true,
		RequiredSecuritySchemes: []string{"BearerAuth"},
	})
	if len(issues) > 0 {
		t.Fatalf("openapi quality violations: %v", issues)
	}

	raw, err := e.MarshalOpenAPIJSON()
	if err != nil {
		t.Fatalf("marshal openapi: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("openapi json is empty")
	}

	requiredOperationIDs := []string{
		"auth_register", "auth_login", "auth_refresh", "auth_me",
		"users_me_get", "users_me_patch",
		"users_me_settings_get", "users_me_settings_put", "users_me_settings_patch",
		"workspaces_create", "workspaces_list", "workspaces_get", "workspaces_patch", "workspaces_delete",
		"workspace_users_list", "workspace_users_add", "workspace_users_patch", "workspace_users_delete",
		"bookmarks_create", "bookmarks_list", "bookmarks_get", "bookmarks_patch", "bookmarks_delete",
		"tags_create", "tags_list", "tags_patch", "tags_delete",
		"admin_users_list", "admin_users_get",
	}

	var openapi map[string]any
	if err := json.Unmarshal(raw, &openapi); err != nil {
		t.Fatalf("decode openapi json: %v", err)
	}

	payload := string(raw)
	for _, id := range requiredOperationIDs {
		if !strings.Contains(payload, `"operationId": "`+id+`"`) {
			t.Fatalf("required operationId %q not found", id)
		}
	}
	components, _ := openapi["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	keyRe := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	for k := range schemas {
		if !keyRe.MatchString(k) {
			t.Fatalf("components.schemas key %q is not OpenAPI-safe", k)
		}
	}
	paths, _ := openapi["paths"].(map[string]any)
	assertResponseStatus(t, paths, "/api/v1/auth/login", "post", "401")
	assertResponseStatus(t, paths, "/api/v1/workspaces/{workspaceId}", "patch", "403")
	assertResponseStatus(t, paths, "/api/v1/workspaces/{workspaceId}/users/{userId}", "patch", "409")
	assertResponseStatus(t, paths, "/api/v1/workspaces/{workspaceId}/bookmarks/{bookmarkId}", "patch", "422")
	assertResponseStatus(t, paths, "/api/v1/workspaces/{workspaceId}/tags/{tagId}", "patch", "404")
}

func assertResponseStatus(t *testing.T, paths map[string]any, path, method, status string) {
	t.Helper()
	pathNode, ok := paths[path].(map[string]any)
	if !ok {
		t.Fatalf("path %s not found in openapi", path)
	}
	opNode, ok := pathNode[method].(map[string]any)
	if !ok {
		t.Fatalf("method %s for path %s not found in openapi", method, path)
	}
	responses, ok := opNode["responses"].(map[string]any)
	if !ok {
		t.Fatalf("responses not found for %s %s", method, path)
	}
	resp, ok := responses[status].(map[string]any)
	if !ok {
		t.Fatalf("response status %s not declared for %s %s", status, method, path)
	}
	content, _ := resp["content"].(map[string]any)
	if _, ok := content["application/problem+json"]; !ok {
		t.Fatalf("response status %s for %s %s has no application/problem+json content", status, method, path)
	}
}
