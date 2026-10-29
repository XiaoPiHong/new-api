package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/jsplugin"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setRelayServiceAuthForRouterTest(t *testing.T, enabled bool, secret string) {
	t.Helper()
	common.OptionMapRWMutex.Lock()
	if common.OptionMap == nil {
		common.OptionMap = make(map[string]string)
	}
	previousEnabled, hadEnabled := common.OptionMap[common.RelayServiceAuthEnabledKey]
	previousSecret, hadSecret := common.OptionMap[common.RelayServiceAuthSecretKey]
	common.OptionMap[common.RelayServiceAuthEnabledKey] = strconv.FormatBool(enabled)
	common.OptionMap[common.RelayServiceAuthSecretKey] = secret
	common.OptionMapRWMutex.Unlock()
	t.Cleanup(func() {
		common.OptionMapRWMutex.Lock()
		defer common.OptionMapRWMutex.Unlock()
		if hadEnabled {
			common.OptionMap[common.RelayServiceAuthEnabledKey] = previousEnabled
		} else {
			delete(common.OptionMap, common.RelayServiceAuthEnabledKey)
		}
		if hadSecret {
			common.OptionMap[common.RelayServiceAuthSecretKey] = previousSecret
		} else {
			delete(common.OptionMap, common.RelayServiceAuthSecretKey)
		}
	})
}

func TestRelayServiceAuthCoversStaticAndNativeRelayRoutes(t *testing.T) {
	setupRelayRouterTestDB(t)
	require.NoError(t, model.DB.AutoMigrate(&model.Midjourney{}))
	static := gin.New()
	SetRelayRouter(static)
	SetTaskPluginProtocolRouter(static)
	SetVideoRouter(static)
	SetTaskRouter(static)
	plugin := compileRouterPlugin(t, "relay-gate", "1.0.0", `[
		{method:"POST",path:"/vendor/gated/jobs",type:"submit",decode:"decode",render:"native"},
		{method:"GET",path:"/vendor/gated/jobs/:task_id",type:"query",render:"native"}
	]`)
	native, registry := newPluginRouterTest(t, []*jsplugin.LoadedPlugin{plugin}, productionPluginRouteHandlers)
	native.NoRoute((&pluginRouteDispatcher{registry: registry}).dispatch)
	for _, route := range []struct {
		engine         http.Handler
		method, path   string
		ordinaryStatus int
	}{
		{static, http.MethodGet, "/v1/models", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/chat/completions", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/responses", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/responses", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/responses/response_missing", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/images/generations", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/images/edits", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/videos", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/videos/task_missing", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/video/generations", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/video/generations/task_missing", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/videos/task_missing/remix", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1beta/models", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1beta/openai/models", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1beta/models/gemini:generateContent", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1beta/models/gemini:countTokens", http.StatusNotFound},
		{static, http.MethodPost, "/pg/chat/completions", http.StatusUnauthorized},
		{static, http.MethodPost, "/mj/submit/imagine", http.StatusUnauthorized},
		{static, http.MethodGet, "/mj/task/task_missing/fetch", http.StatusUnauthorized},
		{static, http.MethodGet, "/mj/image/task_missing", http.StatusUnauthorized},
		{static, http.MethodPost, "/fast/mj/submit/imagine", http.StatusUnauthorized},
		{static, http.MethodGet, "/fast/mj/task/task_missing/fetch", http.StatusUnauthorized},
		{static, http.MethodGet, "/fast/mj/image/task_missing", http.StatusUnauthorized},
		{static, http.MethodPost, "/v1/tasks/relay-gate", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/tasks/task_missing", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/tasks/task_missing/artifacts", http.StatusUnauthorized},
		{static, http.MethodGet, "/v1/tasks/task_missing/artifacts/video/content", http.StatusUnauthorized},
		{static, http.MethodHead, "/v1/tasks/task_missing/artifacts/video/content", http.StatusUnauthorized},
		{native, http.MethodPost, "/vendor/gated/jobs", http.StatusUnauthorized},
		{native, http.MethodGet, "/vendor/gated/jobs/task_missing", http.StatusUnauthorized},
	} {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			for _, state := range []struct {
				name     string
				enabled  bool
				expected string
				actual   []string
				denied   bool
			}{
				{"enabled missing service key", true, "test-relay-secret", nil, true},
				{"enabled wrong service key", true, "test-relay-secret", []string{"incorrect"}, true},
				{"enabled missing configured key", true, "", []string{"test-relay-secret"}, true},
				{"enabled matching service key", true, "test-relay-secret", []string{"test-relay-secret"}, false},
				{"enabled repeated service key", true, "test-relay-secret", []string{"test-relay-secret", "test-relay-secret"}, true},
				{"disabled", false, "test-relay-secret", nil, false},
			} {
				t.Run(state.name, func(t *testing.T) {
					setRelayServiceAuthForRouterTest(t, state.enabled, state.expected)
					request := httptest.NewRequest(route.method, route.path, strings.NewReader(`{}`))
					request.Header.Set("Content-Type", "application/json")
					for _, value := range state.actual {
						request.Header.Add("X-Xphai-Relay-Secret", value)
					}
					recorder := httptest.NewRecorder()
					route.engine.ServeHTTP(recorder, request)
					if state.denied {
						require.Equal(t, http.StatusForbidden, recorder.Code, recorder.Body.String())
						assert.Contains(t, recorder.Body.String(), `"code":"relay_service_access_denied"`)
						assert.NotContains(t, recorder.Body.String(), "test-relay-secret")
						return
					}
					assert.Equal(t, route.ordinaryStatus, recorder.Code, recorder.Body.String())
					assert.NotContains(t, recorder.Body.String(), "relay_service_access_denied")
				})
			}
		})
	}
}

func TestRelayServiceAuthDoesNotInterceptPluginMissFallback(t *testing.T) {
	setRelayServiceAuthForRouterTest(t, true, "test-relay-secret")
	plugin := compileRouterPlugin(t, "relay-gate-miss", "1.0.0", `[{method:"GET",path:"/vendor/gated/:task_id",type:"query",render:"native"}]`)
	outer, registry := newPluginRouterTest(t, []*jsplugin.LoadedPlugin{plugin}, productionPluginRouteHandlers)
	outer.NoRoute((&pluginRouteDispatcher{registry: registry}).dispatch, func(c *gin.Context) { c.String(http.StatusOK, "frontend fallback") })
	recorder := performPluginRequest(outer, http.MethodGet, "/ordinary-frontend-page")
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	assert.Equal(t, "frontend fallback", recorder.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.NotContains(t, recorder.Body.String(), "relay_service_access_denied")
}

func TestMidjourneyImageRequiresAuthenticatedOwnerOrDatabaseAdmin(t *testing.T) {
	setupRelayRouterTestDB(t)
	require.NoError(t, model.DB.AutoMigrate(&model.Midjourney{}, &model.Channel{}))
	previousMemoryCache := common.MemoryCacheEnabled
	common.MemoryCacheEnabled = false
	previousFetch := *system_setting.GetFetchSetting()
	*system_setting.GetFetchSetting() = system_setting.FetchSetting{EnableSSRFProtection: true, AllowPrivateIp: true, AllowedPorts: []string{"1-65535"}}
	t.Cleanup(func() {
		common.MemoryCacheEnabled = previousMemoryCache
		*system_setting.GetFetchSetting() = previousFetch
	})
	service.InitHttpClient()
	var upstreamCalls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		upstreamCalls.Add(1)
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("mock MJ image"))
	}))
	t.Cleanup(upstream.Close)
	users := []model.User{
		{Username: "mj-image-owner", AffCode: "mjowner", Status: common.UserStatusEnabled, Role: common.RoleCommonUser, Group: "default"},
		{Username: "mj-image-other", AffCode: "mjother", Status: common.UserStatusEnabled, Role: common.RoleCommonUser, Group: "default"},
		{Username: "mj-image-admin", AffCode: "mjadmin", Status: common.UserStatusEnabled, Role: common.RoleAdminUser, Group: "default"},
	}
	keys := []string{"mjownerkey", "mjotherkey", "mjadminkey"}
	for i := range users {
		require.NoError(t, model.DB.Create(&users[i]).Error)
		require.NoError(t, model.DB.Create(&model.Token{UserId: users[i].Id, Key: keys[i], Status: common.TokenStatusEnabled, ExpiredTime: -1, UnlimitedQuota: true}).Error)
	}
	channel := &model.Channel{Name: "mj-image-upstream", Type: constant.ChannelTypeMidjourney, Key: "mock-channel-key", BaseURL: &upstream.URL}
	require.NoError(t, model.DB.Create(channel).Error)
	require.NoError(t, model.DB.Create(&model.Midjourney{MjId: "mj_owned_image", UserId: users[0].Id, ChannelId: channel.Id, ImageUrl: upstream.URL + "/image.png"}).Error)
	engine := gin.New()
	SetRelayRouter(engine)
	for _, tc := range []struct {
		name, path, token, serviceKey, roleHeader string
		enabled                                   bool
		status                                    int
	}{
		{name: "owner still needs service key", token: keys[0], enabled: true, status: http.StatusForbidden},
		{name: "service key alone is not user auth", serviceKey: "test-relay-secret", enabled: true, status: http.StatusUnauthorized},
		{name: "owner can read", token: keys[0], serviceKey: "test-relay-secret", enabled: true, status: http.StatusOK},
		{name: "owner can read mode route", path: "/fast/mj/image/mj_owned_image", token: keys[0], serviceKey: "test-relay-secret", enabled: true, status: http.StatusOK},
		{name: "other user cannot read", token: keys[1], serviceKey: "test-relay-secret", enabled: true, status: http.StatusBadRequest},
		{name: "forged admin header cannot read", token: keys[1], serviceKey: "test-relay-secret", roleHeader: strconv.Itoa(common.RoleRootUser), enabled: true, status: http.StatusBadRequest},
		{name: "missing task is indistinguishable", path: "/mj/image/mj_missing_image", token: keys[0], serviceKey: "test-relay-secret", enabled: true, status: http.StatusBadRequest},
		{name: "database admin can read", token: keys[2], serviceKey: "test-relay-secret", enabled: true, status: http.StatusOK},
		{name: "disabled service gate does not allow anonymous", enabled: false, status: http.StatusUnauthorized},
		{name: "disabled service gate still checks ownership", token: keys[1], enabled: false, status: http.StatusBadRequest},
		{name: "disabled service gate permits authenticated owner", token: keys[0], enabled: false, status: http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setRelayServiceAuthForRouterTest(t, tc.enabled, "test-relay-secret")
			path := tc.path
			if path == "" {
				path = "/mj/image/mj_owned_image"
			}
			request := httptest.NewRequest(http.MethodGet, path, nil)
			if tc.token != "" {
				request.Header.Set("Authorization", "Bearer "+tc.token)
			}
			if tc.serviceKey != "" {
				request.Header.Set("X-Xphai-Relay-Secret", tc.serviceKey)
			}
			if tc.roleHeader != "" {
				request.Header.Set("role", tc.roleHeader)
				request.Header.Set("New-Api-User", strconv.Itoa(users[0].Id))
			}
			before := upstreamCalls.Load()
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)
			require.Equal(t, tc.status, recorder.Code, recorder.Body.String())
			if tc.status == http.StatusOK {
				assert.Equal(t, "mock MJ image", recorder.Body.String())
				assert.Equal(t, "image/png", recorder.Header().Get("Content-Type"))
				assert.Equal(t, before+1, upstreamCalls.Load())
				return
			}
			assert.Equal(t, before, upstreamCalls.Load(), "rejected requests must not fetch upstream content")
			assert.NotContains(t, recorder.Body.String(), "mock MJ image")
			if tc.status == http.StatusBadRequest {
				assert.JSONEq(t, `{"error":"midjourney_task_not_found"}`, recorder.Body.String())
			}
		})
	}
}

func TestRelayServiceAuthAllowsOnlyExactSignedArtifactPreview(t *testing.T) {
	setupRelayRouterTestDB(t)
	require.NoError(t, model.DB.AutoMigrate(&model.Task{}))
	previousCryptoSecret, previousMemoryCache := common.CryptoSecret, common.MemoryCacheEnabled
	common.CryptoSecret, common.MemoryCacheEnabled = "test-artifact-signing-secret", false
	t.Cleanup(func() { common.CryptoSecret, common.MemoryCacheEnabled = previousCryptoSecret, previousMemoryCache })
	setRelayServiceAuthForRouterTest(t, true, "test-relay-secret")
	owner := &model.User{Username: "signed-preview-owner", Status: common.UserStatusEnabled, Group: "default"}
	require.NoError(t, model.DB.Create(owner).Error)
	require.NoError(t, model.DB.Create(&model.Token{UserId: owner.Id, Key: "previewtestkey", Status: common.TokenStatusEnabled, ExpiredTime: -1, UnlimitedQuota: true}).Error)
	task := &model.Task{TaskID: "task_signed_preview", UserId: owner.Id, Platform: constant.TaskPlatform("55"), Action: constant.TaskActionTextToVideo, Status: model.TaskStatusSuccess, Progress: "100%", PrivateData: model.TaskPrivateData{ResultURL: "data:video/mp4;base64,cHJldmlldy1ieXRlcw=="}}
	require.NoError(t, model.DB.Create(task).Error)
	access, err := service.IssueTaskArtifactAccess(task.TaskID, "video")
	require.NoError(t, err)
	require.True(t, service.VerifyTaskArtifactAccess(access, task.TaskID, "video"))
	query := "?access=" + url.QueryEscape(access)
	contentPath := "/v1/tasks/" + task.TaskID + "/artifacts/video/content"
	engine := gin.New()
	SetRelayRouter(engine)
	SetTaskPluginProtocolRouter(engine)
	SetVideoRouter(engine)
	SetTaskRouter(engine)
	for _, tc := range []struct {
		name, method, path string
		status             int
		body               string
	}{
		{"signed exact GET", http.MethodGet, contentPath + query, http.StatusOK, "preview-bytes"},
		{"signed exact HEAD", http.MethodHead, contentPath + query, http.StatusOK, ""},
		{"signature cannot authorize other artifact", http.MethodGet, "/v1/tasks/" + task.TaskID + "/artifacts/other/content" + query, http.StatusNotFound, ""},
		{"signature cannot authorize other task", http.MethodGet, "/v1/tasks/task_other/artifacts/video/content" + query, http.StatusNotFound, ""},
		{"tampered signature", http.MethodGet, contentPath + "?access=" + url.QueryEscape(access+"A"), http.StatusNotFound, ""},
		{"unsigned content needs service key", http.MethodGet, contentPath, http.StatusForbidden, ""},
		{"signature cannot submit", http.MethodPost, "/v1/tasks/sora" + query, http.StatusForbidden, ""},
		{"signature cannot query task", http.MethodGet, "/v1/tasks/" + task.TaskID + query, http.StatusForbidden, ""},
		{"signature cannot list artifacts", http.MethodGet, "/v1/tasks/" + task.TaskID + "/artifacts" + query, http.StatusForbidden, ""},
		{"signature cannot generate video", http.MethodPost, "/v1/videos" + query, http.StatusForbidden, ""},
		{"signature cannot use playground", http.MethodPost, "/pg/chat/completions" + query, http.StatusForbidden, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{}`))
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)
			require.Equal(t, tc.status, recorder.Code, recorder.Body.String())
			if tc.status == http.StatusOK {
				assert.Equal(t, tc.body, recorder.Body.String())
				assert.Equal(t, "video/mp4", recorder.Header().Get("Content-Type"))
			} else {
				assert.NotContains(t, recorder.Body.String(), "preview-bytes")
			}
		})
	}
	for _, state := range []struct {
		name, serviceKey string
		status           int
	}{
		{"API token without service key", "", http.StatusForbidden},
		{"API token with service key", "test-relay-secret", http.StatusOK},
	} {
		t.Run(state.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, contentPath, nil)
			request.Header.Set("Authorization", "Bearer previewtestkey")
			if state.serviceKey != "" {
				request.Header.Set("X-Xphai-Relay-Secret", state.serviceKey)
			}
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)
			require.Equal(t, state.status, recorder.Code, recorder.Body.String())
			if state.status == http.StatusOK {
				assert.Equal(t, "preview-bytes", recorder.Body.String())
			}
		})
	}
	require.NoError(t, model.DB.Delete(task).Error)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, contentPath+query, nil))
	assert.Equal(t, http.StatusNotFound, recorder.Code, "a signature must not resurrect a deleted resource")
}

func TestListModelsSupportsOpenAIAndGeminiAuthentication(t *testing.T) {
	setupRelayRouterTestDB(t)

	user := model.User{
		Username: "models-user",
		Status:   common.UserStatusEnabled,
		Group:    "default",
		Quota:    100,
	}
	require.NoError(t, model.DB.Create(&user).Error)
	require.NoError(t, model.DB.Create(&model.Token{
		UserId:         user.Id,
		Key:            "modelstestkey",
		Status:         common.TokenStatusEnabled,
		ExpiredTime:    -1,
		UnlimitedQuota: true,
	}).Error)

	engine := gin.New()
	SetRelayRouter(engine)

	tests := []struct {
		name           string
		path           string
		headerName     string
		expectedObject string
		expectedField  string
	}{
		{
			name:           "OpenAI bearer token",
			path:           "/v1/models",
			headerName:     "Authorization",
			expectedObject: "list",
			expectedField:  "data",
		},
		{
			name:          "Gemini API key header",
			path:          "/v1/models",
			headerName:    "x-goog-api-key",
			expectedField: "models",
		},
		{
			name:          "Gemini API key query",
			path:          "/v1/models?key=modelstestkey",
			expectedField: "models",
		},
		{
			name:          "Gemini native model list",
			path:          "/v1beta/models",
			headerName:    "x-goog-api-key",
			expectedField: "models",
		},
		{
			name:           "Gemini compatible model list",
			path:           "/v1beta/openai/models",
			headerName:     "x-goog-api-key",
			expectedObject: "list",
			expectedField:  "data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, state := range []struct {
				name           string
				enabled        bool
				header         string
				expectedStatus int
			}{
				{"disabled", false, "", http.StatusOK},
				{"enabled missing service key", true, "", http.StatusForbidden},
				{"enabled wrong service key", true, "incorrect", http.StatusForbidden},
				{"enabled matching service key", true, "test-relay-secret", http.StatusOK},
			} {
				t.Run(state.name, func(t *testing.T) {
					setRelayServiceAuthForRouterTest(t, state.enabled, "test-relay-secret")
					recorder := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodGet, test.path, nil)
					if test.headerName != "" {
						value := "modelstestkey"
						if test.headerName == "Authorization" {
							value = "Bearer " + value
						}
						request.Header.Set(test.headerName, value)
					}
					if state.header != "" {
						request.Header.Set("X-Xphai-Relay-Secret", state.header)
					}
					engine.ServeHTTP(recorder, request)
					require.Equal(t, state.expectedStatus, recorder.Code, recorder.Body.String())
					if state.expectedStatus == http.StatusForbidden {
						assert.Contains(t, recorder.Body.String(), "relay_service_access_denied")
						return
					}
					var payload map[string]any
					require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &payload))
					assert.Contains(t, payload, test.expectedField)
					assert.NotContains(t, payload, "error")
					if test.expectedObject != "" {
						assert.Equal(t, test.expectedObject, payload["object"])
					}
				})
			}
		})
	}
}

func setupRelayRouterTestDB(t *testing.T) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	originalIsMasterNode := common.IsMasterNode
	originalRedisEnabled := common.RedisEnabled
	originalSQLitePath := common.SQLitePath
	originalMainDatabaseType := common.MainDatabaseType()
	originalLogDatabaseType := common.LogDatabaseType()
	originalSQLDSN, hadSQLDSN := os.LookupEnv("SQL_DSN")

	common.IsMasterNode = false
	common.RedisEnabled = false
	common.SQLitePath = fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	common.SetDatabaseTypes(common.DatabaseTypeSQLite, common.DatabaseTypeSQLite)
	require.NoError(t, os.Setenv("SQL_DSN", "local"))
	database, err := gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{})
	require.NoError(t, err)
	model.DB = database
	require.NoError(t, model.DB.AutoMigrate(&model.User{}, &model.Token{}, &model.Ability{}, &model.Option{}))
	// Initialize the real dialect helpers after preparing the minimal schema.
	// A non-master InitDB must be able to read the legacy retirement option.
	require.NoError(t, model.InitDB())
	bootstrapDB, err := database.DB()
	require.NoError(t, err)
	require.NoError(t, bootstrapDB.Close())
	model.LOG_DB = model.DB

	t.Cleanup(func() {
		if sqlDB, err := model.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
		common.IsMasterNode = originalIsMasterNode
		common.RedisEnabled = originalRedisEnabled
		common.SQLitePath = originalSQLitePath
		common.SetDatabaseTypes(originalMainDatabaseType, originalLogDatabaseType)
		if hadSQLDSN {
			require.NoError(t, os.Setenv("SQL_DSN", originalSQLDSN))
		} else {
			require.NoError(t, os.Unsetenv("SQL_DSN"))
		}
	})
}
