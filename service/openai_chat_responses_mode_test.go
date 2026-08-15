package service

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/stretchr/testify/require"
)

func TestShouldPerplexityAgentModelUseResponses(t *testing.T) {
	tests := []struct {
		name        string
		channelType int
		model       string
		want        bool
	}{
		{name: "anthropic agent model", channelType: constant.ChannelTypePerplexity, model: "anthropic/claude-opus-4-8", want: true},
		{name: "openai agent model", channelType: constant.ChannelTypePerplexity, model: "openai/gpt-5.6-sol", want: true},
		{name: "perplexity agent model", channelType: constant.ChannelTypePerplexity, model: "perplexity/sonar", want: true},
		{name: "legacy perplexity chat model", channelType: constant.ChannelTypePerplexity, model: "sonar", want: false},
		{name: "other channel", channelType: constant.ChannelTypeOpenAI, model: "anthropic/claude-opus-4-8", want: false},
		{name: "trimmed model", channelType: constant.ChannelTypePerplexity, model: "  google/gemini-3.5-flash  ", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ShouldPerplexityAgentModelUseResponses(tt.channelType, tt.model))
		})
	}
}
