package graph

import (
	chatmodel "audio2markdown/component/chat_model"
	fileio "audio2markdown/component/file_io"
	"audio2markdown/component/logic"
	"audio2markdown/component/prompt"
	"audio2markdown/component/retriever"
	"context"

	"github.com/cloudwego/eino/compose"
)

// BuildAgentGraph 构建面试总结 Agent Graph
func BuildAgentGraph(ctx context.Context) (compose.Runnable[string, string], error) {
	const (
		ReadInputNode   = "ReadInput"
		MessagesToQuery = "MessagesToQuery"
		RetrieveNode    = "Retriever"
		PromptNode      = "Prompt"
		ChatModelNode   = "ChatModel"
		WriteOutputNode = "WriteOutput"
	)

	g := compose.NewGraph[string, string]()

	_ = g.AddLambdaNode(ReadInputNode, compose.InvokableLambda(fileio.NewReadConversationLambda))
	_ = g.AddLambdaNode(MessagesToQuery, compose.InvokableLambda(logic.MessagesToQueryLambda))

	retrieverNode, err := retriever.NewRedisRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RetrieveNode, retrieverNode, compose.WithOutputKey("knowledge_docs"))

	chatTpl := prompt.NewInterviewSummaryPrompt()
	_ = g.AddChatTemplateNode(PromptNode, chatTpl)

	chatModel := chatmodel.NewOpenAIChatModel(ctx)
	_ = g.AddChatModelNode(ChatModelNode, chatModel)

	_ = g.AddLambdaNode(WriteOutputNode, compose.InvokableLambda(fileio.NewWriteMarkdownLambda))

	// Edges
	_ = g.AddEdge(compose.START, ReadInputNode)
	_ = g.AddEdge(ReadInputNode, MessagesToQuery)
	_ = g.AddEdge(MessagesToQuery, RetrieveNode)
	_ = g.AddEdge(MessagesToQuery, PromptNode)
	_ = g.AddEdge(RetrieveNode, PromptNode)
	_ = g.AddEdge(PromptNode, ChatModelNode)
	_ = g.AddEdge(ChatModelNode, WriteOutputNode)
	_ = g.AddEdge(WriteOutputNode, compose.END)

	r, err := g.Compile(ctx, compose.WithGraphName("InterviewAgentGraph"))
	if err != nil {
		return nil, err
	}
	return r, nil
}
