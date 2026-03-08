package bms_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/app"
	"github.com/pafthang/bms/internal/config"
	dbpkg "github.com/pafthang/bms/internal/db"
	"github.com/pafthang/bms/internal/domain/models"
	"github.com/pafthang/dbx"
	"github.com/pafthang/orm"
)

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authRegisterResp struct {
	Data struct {
		User struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Tokens tokenPair `json:"tokens"`
	} `json:"data"`
}

type workspaceCreateResp struct {
	Data struct {
		Workspace struct {
			ID int64 `json:"id"`
		} `json:"workspace"`
	} `json:"data"`
}

type workspaceMembersResp struct {
	Data []struct {
		UserID int64  `json:"user_id"`
		Role   string `json:"role"`
	} `json:"data"`
}

type workspacesListResp struct {
	Data []struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

type apiErrorResp struct {
	Code string `json:"code"`
}

type bookmarkCreateResp struct {
	Data struct {
		ID     int64   `json:"id"`
		TagIDs []int64 `json:"tag_ids"`
	} `json:"data"`
}

type testEnv struct {
	t      *testing.T
	db     *dbx.DB
	engine *arc.Engine
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "bms-test.sqlite")
	cfg := config.Config{
		HTTPAddr:        ":0",
		ShutdownTimeout: 3 * time.Second,
		AllowedOrigins:  []string{"http://localhost:3000"},
		DBDSN:           fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", dbPath),
		JWTSecret:       "test-secret",
		AccessTTL:       time.Hour,
		RefreshTTL:      24 * time.Hour,
	}
	database, err := dbpkg.Open(cfg)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if err := dbpkg.MigrateUp(context.Background(), database, filepath.Join("migrations")); err != nil {
		t.Fatalf("migrate up: %v", err)
	}
	engine := app.BuildEngine(cfg, database, app.BuildOptions{IncludeSystemRoutes: true, IncludeHealthRoutes: true})
	return &testEnv{t: t, db: database, engine: engine}
}

func (e *testEnv) doJSON(method, path string, token string, body any, out any) (int, []byte) {
	e.t.Helper()
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			e.t.Fatalf("marshal request: %v", err)
		}
	} else {
		reqBody = []byte("{}")
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)
	if out != nil {
		if err := json.Unmarshal(w.Body.Bytes(), out); err != nil {
			e.t.Fatalf("unmarshal response (%s %s): %v; body=%s", method, path, err, w.Body.String())
		}
	}
	return w.Code, w.Body.Bytes()
}

func (e *testEnv) register(email string) (int64, tokenPair) {
	e.t.Helper()
	var resp authRegisterResp
	status, body := e.doJSON(http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"email":    email,
		"password": "password123",
	}, &resp)
	if status != http.StatusCreated {
		e.t.Fatalf("register status=%d body=%s", status, string(body))
	}
	return resp.Data.User.ID, resp.Data.Tokens
}

func (e *testEnv) login(email string) tokenPair {
	e.t.Helper()
	var resp authRegisterResp
	status, body := e.doJSON(http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": "password123",
	}, &resp)
	if status != http.StatusOK {
		e.t.Fatalf("login status=%d body=%s", status, string(body))
	}
	return resp.Data.Tokens
}

func TestAcceptanceCoreFlows(t *testing.T) {
	env := setupTestEnv(t)

	ownerID, ownerTokens := env.register("owner@example.com")
	viewerID, viewerTokens := env.register("viewer@example.com")

	var meResp map[string]any
	status, body := env.doJSON(http.MethodGet, "/api/v1/auth/me", ownerTokens.AccessToken, nil, &meResp)
	if status != http.StatusOK {
		t.Fatalf("auth/me status=%d body=%s", status, string(body))
	}

	var wsCreate workspaceCreateResp
	status, body = env.doJSON(http.MethodPost, "/api/v1/workspaces", ownerTokens.AccessToken, map[string]any{
		"name":        "Main Workspace",
		"description": "demo",
	}, &wsCreate)
	if status != http.StatusCreated {
		t.Fatalf("create workspace status=%d body=%s", status, string(body))
	}
	workspaceID := wsCreate.Data.Workspace.ID

	var members workspaceMembersResp
	status, body = env.doJSON(http.MethodGet, fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), ownerTokens.AccessToken, nil, &members)
	if status != http.StatusOK {
		t.Fatalf("list members status=%d body=%s", status, string(body))
	}
	if len(members.Data) == 0 || members.Data[0].UserID != ownerID || members.Data[0].Role != "admin" {
		t.Fatalf("creator must be admin: %+v", members.Data)
	}

	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), ownerTokens.AccessToken, map[string]any{
		"user_id": viewerID,
		"role":    "viewer",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("add viewer status=%d body=%s", status, string(body))
	}

	var workspaceList workspacesListResp
	status, body = env.doJSON(http.MethodGet, "/api/v1/workspaces", viewerTokens.AccessToken, nil, &workspaceList)
	if status != http.StatusOK || len(workspaceList.Data) == 0 {
		t.Fatalf("viewer should see workspace: status=%d body=%s", status, string(body))
	}

	var errResp apiErrorResp
	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/bookmarks", workspaceID), viewerTokens.AccessToken, map[string]any{
		"title": "Denied",
		"url":   "https://example.com",
	}, &errResp)
	if status != http.StatusForbidden || errResp.Code != "workspace_forbidden" {
		t.Fatalf("viewer create bookmark should be forbidden: status=%d code=%s body=%s", status, errResp.Code, string(body))
	}

	status, body = env.doJSON(http.MethodPatch, fmt.Sprintf("/api/v1/workspaces/%d/users/%d", workspaceID, viewerID), ownerTokens.AccessToken, map[string]any{
		"role": "editor",
	}, nil)
	if status != http.StatusOK {
		t.Fatalf("promote to editor status=%d body=%s", status, string(body))
	}

	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/bookmarks", workspaceID), viewerTokens.AccessToken, map[string]any{
		"title": "Allowed",
		"url":   "https://example.com",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("editor create bookmark should pass: status=%d body=%s", status, string(body))
	}

	superID, _ := env.register("super@example.com")
	superUser, err := orm.ByPK[models.User](context.Background(), env.db, superID)
	if err != nil {
		t.Fatalf("load super user: %v", err)
	}
	superUser.IsSuperadmin = true
	if err := orm.Update(context.Background(), env.db, superUser); err != nil {
		t.Fatalf("promote superadmin: %v", err)
	}
	superTokens := env.login("super@example.com")
	status, body = env.doJSON(http.MethodGet, "/api/v1/auth/me", superTokens.AccessToken, nil, &meResp)
	if status != http.StatusOK || meResp["data"].(map[string]any)["is_superadmin"] != true {
		t.Fatalf("superadmin token should carry super flag: status=%d body=%s", status, string(body))
	}

	status, body = env.doJSON(http.MethodGet, "/api/v1/workspaces?all=true", superTokens.AccessToken, nil, &workspaceList)
	if status != http.StatusOK || len(workspaceList.Data) == 0 {
		t.Fatalf("superadmin should see all workspaces: status=%d body=%s", status, string(body))
	}

	status, body = env.doJSON(http.MethodGet, fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), superTokens.AccessToken, nil, &members)
	if status != http.StatusOK || len(members.Data) < 2 {
		t.Fatalf("superadmin should list members: status=%d body=%s", status, string(body))
	}
}

func TestMembershipConflict(t *testing.T) {
	env := setupTestEnv(t)
	_, ownerTokens := env.register("owner2@example.com")
	memberID, _ := env.register("member2@example.com")

	var wsCreate workspaceCreateResp
	status, body := env.doJSON(http.MethodPost, "/api/v1/workspaces", ownerTokens.AccessToken, map[string]any{
		"name": "Conflict Workspace",
	}, &wsCreate)
	if status != http.StatusCreated {
		t.Fatalf("create workspace status=%d body=%s", status, string(body))
	}
	workspaceID := wsCreate.Data.Workspace.ID

	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), ownerTokens.AccessToken, map[string]any{
		"user_id": memberID,
		"role":    "viewer",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("first add member status=%d body=%s", status, string(body))
	}

	var errResp apiErrorResp
	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/users", workspaceID), ownerTokens.AccessToken, map[string]any{
		"user_id": memberID,
		"role":    "viewer",
	}, &errResp)
	if status != http.StatusConflict || errResp.Code != "workspace_membership_conflict" {
		t.Fatalf("duplicate member should conflict: status=%d code=%s body=%s", status, errResp.Code, string(body))
	}
}

func TestLastAdminGuard(t *testing.T) {
	env := setupTestEnv(t)
	ownerID, ownerTokens := env.register("owner3@example.com")

	var wsCreate workspaceCreateResp
	status, body := env.doJSON(http.MethodPost, "/api/v1/workspaces", ownerTokens.AccessToken, map[string]any{
		"name": "Last Admin Workspace",
	}, &wsCreate)
	if status != http.StatusCreated {
		t.Fatalf("create workspace status=%d body=%s", status, string(body))
	}
	workspaceID := wsCreate.Data.Workspace.ID

	var errResp apiErrorResp
	status, body = env.doJSON(http.MethodPatch, fmt.Sprintf("/api/v1/workspaces/%d/users/%d", workspaceID, ownerID), ownerTokens.AccessToken, map[string]any{
		"role": "viewer",
	}, &errResp)
	if status != http.StatusConflict || errResp.Code != "workspace_last_admin_violation" {
		t.Fatalf("demote last admin should conflict: status=%d code=%s body=%s", status, errResp.Code, string(body))
	}

	status, body = env.doJSON(http.MethodDelete, fmt.Sprintf("/api/v1/workspaces/%d/users/%d", workspaceID, ownerID), ownerTokens.AccessToken, nil, &errResp)
	if status != http.StatusConflict || errResp.Code != "workspace_last_admin_violation" {
		t.Fatalf("delete last admin should conflict: status=%d code=%s body=%s", status, errResp.Code, string(body))
	}
}

func TestSoftDeletedTagIsIgnoredInBookmarkView(t *testing.T) {
	env := setupTestEnv(t)
	_, ownerTokens := env.register("owner4@example.com")

	var wsCreate workspaceCreateResp
	status, body := env.doJSON(http.MethodPost, "/api/v1/workspaces", ownerTokens.AccessToken, map[string]any{
		"name": "Tags Workspace",
	}, &wsCreate)
	if status != http.StatusCreated {
		t.Fatalf("create workspace status=%d body=%s", status, string(body))
	}
	workspaceID := wsCreate.Data.Workspace.ID

	var tagResp struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/tags", workspaceID), ownerTokens.AccessToken, map[string]any{
		"name":  "important",
		"color": "#abc",
	}, &tagResp)
	if status != http.StatusCreated {
		t.Fatalf("create tag status=%d body=%s", status, string(body))
	}

	var bookmarkResp bookmarkCreateResp
	status, body = env.doJSON(http.MethodPost, fmt.Sprintf("/api/v1/workspaces/%d/bookmarks", workspaceID), ownerTokens.AccessToken, map[string]any{
		"title":   "Item",
		"url":     "https://example.com",
		"tag_ids": []int64{tagResp.Data.ID},
	}, &bookmarkResp)
	if status != http.StatusCreated || len(bookmarkResp.Data.TagIDs) != 1 {
		t.Fatalf("create bookmark with tag status=%d body=%s", status, string(body))
	}
	bookmarkID := bookmarkResp.Data.ID

	status, body = env.doJSON(http.MethodDelete, fmt.Sprintf("/api/v1/workspaces/%d/tags/%d", workspaceID, tagResp.Data.ID), ownerTokens.AccessToken, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("delete tag status=%d body=%s", status, string(body))
	}

	var bookmarkGet bookmarkCreateResp
	status, body = env.doJSON(http.MethodGet, fmt.Sprintf("/api/v1/workspaces/%d/bookmarks/%d", workspaceID, bookmarkID), ownerTokens.AccessToken, nil, &bookmarkGet)
	if status != http.StatusOK {
		t.Fatalf("get bookmark status=%d body=%s", status, string(body))
	}
	if len(bookmarkGet.Data.TagIDs) != 0 {
		t.Fatalf("soft deleted tag should be ignored in bookmark view, got=%v body=%s", bookmarkGet.Data.TagIDs, string(body))
	}
}
