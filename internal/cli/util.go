package cli

import (
	"context"
	"github.com/anthropics/anthropic-sdk-go"
)

// general purpose util function to get certain inferences via system  prompts from LLMs
func getAnswerByPrompt(input, prompt string) (string, error) {
	client := anthropic.NewClient()
	message := anthropic.NewUserMessage(anthropic.NewTextBlock(input))
	answer, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		System:    []anthropic.TextBlockParam{{Text: prompt}},
		Model:     anthropic.ModelClaude4Sonnet20250514,
		Messages:  []anthropic.MessageParam{message},
		MaxTokens: int64(1024),
	})

	if err != nil {
		return "", err
	}

	var out string
	for _, content := range answer.Content {
		switch content.Type {
		case "text":
			out = content.Text
		}
	}

	return out, nil
}
