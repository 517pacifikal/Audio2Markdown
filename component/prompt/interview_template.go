package prompt

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// NewInterviewSummaryPrompt 返回面试总结专用 Prompt 模板
func NewInterviewSummaryPrompt() prompt.ChatTemplate {
	systemTpl := `你是一名专业的面试总结助手。请根据给定的面试对话内容和知识库，完成以下任务：
1. 总结面试中出现了哪些问题（问题列表）。
2. 针对每个问题，分析候选人的回答与预期的符合程度。
3. 总结候选人的表现中有哪些不足，并指出具体与哪些知识点相关。
请以结构化 Markdown 格式输出，内容包括：问题列表、每个问题的分析、候选人不足及关联知识点。`

	return prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		schema.UserMessage("面试对话内容：{result}\n知识库参考：{knowledge_docs}"),
	)
}
