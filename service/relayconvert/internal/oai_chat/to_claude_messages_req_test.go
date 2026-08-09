package oaichat

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIChatRequestToClaudeMessagesOmitsEmptyTools(t *testing.T) {
	request := dto.GeneralOpenAIRequest{
		Model: "claude-opus-5",
		Messages: []dto.Message{
			{Role: "user", Content: "hello"},
		},
	}

	converted, err := OpenAIChatRequestToClaudeMessages(nil, request)
	require.NoError(t, err)
	assert.Nil(t, converted.Tools)

	encoded, err := common.Marshal(converted)
	require.NoError(t, err)
	assert.NotContains(t, string(encoded), `"tools"`)
}

