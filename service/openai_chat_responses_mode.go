package service

import (
	"strings"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/service/relayconvert"
	"github.com/QuantumNous/new-api/setting/model_setting"
)

// ShouldPerplexityAgentModelUseResponses identifies provider-qualified models
// exposed by Perplexity's Agent API. Those models are available through the
// Responses endpoint, while Perplexity's legacy Chat Completions endpoint only
// accepts its chat model identifiers.
func ShouldPerplexityAgentModelUseResponses(channelType int, model string) bool {
	if channelType != constant.ChannelTypePerplexity {
		return false
	}
	return strings.Contains(strings.TrimSpace(model), "/")
}

func ShouldChatCompletionsUseResponsesPolicy(policy model_setting.ChatCompletionsToResponsesPolicy, channelID int, channelType int, model string) bool {
	return relayconvert.ShouldChatCompletionsUseResponsesPolicy(policy, channelID, channelType, model)
}

func ShouldChatCompletionsUseResponsesGlobal(channelID int, channelType int, model string) bool {
	if ShouldPerplexityAgentModelUseResponses(channelType, model) {
		return true
	}
	return relayconvert.ShouldChatCompletionsUseResponsesGlobal(channelID, channelType, model)
}
