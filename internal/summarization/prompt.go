package summarization

import (
	"fmt"
	"strings"
)

func BuildSummarizationPrompt(text string, minWords, maxWords int, weakChars []string) (system string, user string) {
	system = fmt.Sprintf(`You are a text summarizer for a typing practice application.

Your job is to convert source content into a concise, well-written passage that is pleasant to type.

Rules:
- Output a coherent passage of %d to %d words
- Preserve important concepts and technical terminology
- Preserve the original meaning of the source
- Remove repetition, unnecessary examples, advertisements, and calls to action
- Write natural English prose that flows well
- Do NOT use markdown formatting
- Do NOT use bullet points or numbered lists
- Do NOT include URLs or links
- Do NOT include code blocks or inline code
- Do NOT use special characters unnecessarily
- Avoid excessive use of parentheses, brackets, or uncommon punctuation
- The passage should read like a well-written article paragraph or two
- Focus on the key ideas — this is not a one-sentence summary, but a condensed version of the source`, minWords, maxWords)

	if len(weakChars) > 0 {
		system += fmt.Sprintf(`
- The user struggles with these characters: [%s]. Try to naturally include more words that use these characters, without making the passage feel forced or unnatural`, strings.Join(weakChars, ", "))
	}

	user = fmt.Sprintf("Summarize the following text into a typing passage:\n\n%s", text)
	return system, user
}
