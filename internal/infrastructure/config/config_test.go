package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dev-shimada/gostubby/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	// Ensure directory exists
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	path := filepath.Join(dir, name)
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return path
}

func TestConfigRepository_Load(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 有効なJSONファイルを作成
	validJSON := `[
		{
			"name": "test1",
			"description": "test endpoint 1",
			"request": {
				"url": "/test1",
				"method": "GET",
				"queryParameters": {
					"param1": {
						"equalTo": "value1"
					}
				}
			},
			"response": {
				"status": 200,
				"body": "test response 1"
			}
		}
	]`
	validPath := createTestFile(t, tmpDir, "valid.json", validJSON)

	// サブディレクトリを作成
	subDir := filepath.Join(tmpDir, "subdir")

	// 複数のエンドポイントを含むJSONファイルを作成
	multipleJSON1 := `[
		{
			"name": "test2",
			"description": "test endpoint 2",
			"request": {
				"url": "/test2",
				"method": "POST"
			},
			"response": {
				"status": 201,
				"body": "test response 2"
			}
		}
	]`
	createTestFile(t, subDir, "multiple.json", multipleJSON1)

	// 不正なJSONファイルを作成
	invalidJSON := `{ invalid json`
	invalidPath := createTestFile(t, tmpDir, "invalid.json", invalidJSON)

	// 非JSONファイルを作成
	nonJSONPath := createTestFile(t, tmpDir, "test.txt", "text file")

	repo := NewConfigRepository()

	tests := []struct {
		name    string
		path    string
		want    int // expected number of endpoints
		wantErr bool
	}{
		{
			name:    "valid single JSON file",
			path:    validPath,
			want:    1,
			wantErr: false,
		},
		{
			name:    "directory with multiple JSON files",
			path:    subDir,
			want:    1,
			wantErr: false,
		},
		{
			name:    "invalid JSON file",
			path:    invalidPath,
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-JSON file",
			path:    nonJSONPath,
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-existent path",
			path:    filepath.Join(tmpDir, "nonexistent"),
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoints, err := repo.Load(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, endpoints, tt.want)

			// 追加の検証
			if tt.name == "valid single JSON file" {
				assert.Equal(t, "test1", endpoints[0].Name)
				assert.Equal(t, "/test1", endpoints[0].Request.URL)
				assert.Equal(t, "GET", endpoints[0].Request.Method)
				assert.Equal(t, 200, endpoints[0].Response.Status)
			}
		})
	}
}

func TestEndpointValidation(t *testing.T) {
	validEndpoint := model.Endpoint{
		Name:        "test",
		Description: "test endpoint",
		Request: model.Request{
			URL:    "/test",
			Method: "GET",
			QueryParameters: map[string]model.Matcher{
				"param1": {
					EqualTo: "value1",
				},
			},
		},
		Response: model.Response{
			Status: 200,
			Body:   "test response",
		},
	}

	tmpDir := t.TempDir()
	endpointJSON, err := os.CreateTemp(tmpDir, "endpoint*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := endpointJSON.Close(); err != nil {
			t.Errorf("failed to close file: %v", err)
		}
		if err := os.Remove(endpointJSON.Name()); err != nil {
			t.Errorf("failed to remove file: %v", err)
		}
	}()

	repo := NewConfigRepository()

	// 有効なエンドポイントをJSONとして保存
	validJSON := `[
		{
			"name": "test",
			"description": "test endpoint",
			"request": {
				"url": "/test",
				"method": "GET",
				"queryParameters": {
					"param1": {
						"equalTo": "value1"
					}
				}
			},
			"response": {
				"status": 200,
				"body": "test response"
			}
		}
	]`
	err = os.WriteFile(endpointJSON.Name(), []byte(validJSON), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// 有効なJSONファイルを読み込んでバリデーション
	endpoints, err := repo.Load(endpointJSON.Name())
	assert.NoError(t, err)
	assert.Len(t, endpoints, 1)
	assert.Equal(t, validEndpoint.Name, endpoints[0].Name)
	assert.Equal(t, validEndpoint.Description, endpoints[0].Description)
	assert.Equal(t, validEndpoint.Request.URL, endpoints[0].Request.URL)
	assert.Equal(t, validEndpoint.Request.Method, endpoints[0].Request.Method)
	assert.Equal(t, validEndpoint.Response.Status, endpoints[0].Response.Status)
	assert.Equal(t, validEndpoint.Response.Body, endpoints[0].Response.Body)
}
