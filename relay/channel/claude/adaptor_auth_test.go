package claude

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClaudeAuthTestContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	return c
}

func TestSetupRequestHeaderUsesOAuthBearerForOAT(t *testing.T) {
	c := newClaudeAuthTestContext()
	c.Request.Header.Set("anthropic-beta", "prompt-caching-2024-07-31")
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey:          "  sk-ant-oat01-test  ",
			OriginModelName: "claude-opus-4-7",
		},
	}
	header := http.Header{}

	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.Equal(t, "Bearer sk-ant-oat01-test", header.Get("Authorization"))
	assert.Empty(t, header.Get("x-api-key"))
	assert.Equal(t, "2023-06-01", header.Get("anthropic-version"))
	assert.Contains(t, strings.Split(header.Get("anthropic-beta"), ","), claudeOAuthBeta)
	assert.Contains(t, strings.Split(header.Get("anthropic-beta"), ","), "prompt-caching-2024-07-31")
}

func TestSetupRequestHeaderKeepsAPIKeyAuthentication(t *testing.T) {
	c := newClaudeAuthTestContext()
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey:          "sk-ant-api03-test",
			OriginModelName: "claude-opus-4-7",
		},
	}
	header := http.Header{"Authorization": []string{"Bearer stale-token"}}

	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.Equal(t, "sk-ant-api03-test", header.Get("x-api-key"))
	assert.Empty(t, header.Get("Authorization"))
	assert.NotContains(t, strings.Split(header.Get("anthropic-beta"), ","), claudeOAuthBeta)
}

func TestSetupRequestHeaderForwardsChannelProxyToClaudeAgentGateway(t *testing.T) {
	c := newClaudeAuthTestContext()
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey:          "sk-ant-oat01-test",
			OriginModelName: "claude-opus-5",
			ChannelBaseUrl:  "https://claude-agent-gateway.example.workers.dev",
			ChannelSetting: dto.ChannelSettings{
				Proxy: "http://proxy-user:proxy-pass@proxy.example:8080",
			},
		},
	}
	header := http.Header{}

	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.Equal(t, info.ChannelSetting.Proxy, header.Get(claudeAgentProxyURLHeader))
}

func TestSetupRequestHeaderDoesNotLeakChannelProxyToAnthropic(t *testing.T) {
	c := newClaudeAuthTestContext()
	info := &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiKey:          "sk-ant-api03-test",
			OriginModelName: "claude-opus-4-7",
			ChannelBaseUrl:  "https://api.anthropic.com",
			ChannelSetting: dto.ChannelSettings{
				Proxy: "http://proxy-user:proxy-pass@proxy.example:8080",
			},
		},
	}
	header := http.Header{}

	require.NoError(t, (&Adaptor{}).SetupRequestHeader(c, &header, info))
	assert.Empty(t, header.Get(claudeAgentProxyURLHeader))
}
