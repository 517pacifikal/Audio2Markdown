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
		RetrieveNode    = "Retriever"
		PromptNode      = "Prompt"
		ChatModelNode   = "ChatModel"
		WriteOutputNode = "WriteOutput"
		MergeInputNode  = "MergeInput"
	)

	g := compose.NewGraph[string, string]()

	// Lambda节点：读取目录下txt
	_ = g.AddLambdaNode(ReadInputNode, compose.InvokableLambda(fileio.NewReadConversationLambda))

	// 检索节点
	retrieverNode, err := retriever.NewRedisRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(RetrieveNode, retrieverNode, compose.WithOutputKey("knowledge_docs"))

	// ChatTemplate节点
	chatTpl := prompt.NewInterviewSummaryPrompt()
	_ = g.AddChatTemplateNode(PromptNode, chatTpl)

	// ChatModel节点
	chatModel := chatmodel.NewOpenAIChatModel(ctx)
	_ = g.AddChatModelNode(ChatModelNode, chatModel)

	// Lambda节点：写入md
	_ = g.AddLambdaNode(WriteOutputNode, compose.InvokableLambda(fileio.NewWriteMarkdownLambda))

	_ = g.AddLambdaNode(MergeInputNode, compose.InvokableLambda(logic.NewMergeInputNode))

	// Edges
	_ = g.AddEdge(compose.START, ReadInputNode)
	_ = g.AddEdge(ReadInputNode, RetrieveNode)    // ReadInputNode -> RetrieveNode
	_ = g.AddEdge(RetrieveNode, MergeInputNode)   // RetrieveNode -> MergeInputNode
	_ = g.AddEdge(ReadInputNode, MergeInputNode)  // ReadInputNode -> MergeInputNode
	_ = g.AddEdge(RetrieveNode, MergeInputNode)   // RetrieveNode -> MergeInputNode
	_ = g.AddEdge(MergeInputNode, PromptNode)     // MergeInputNode -> PromptNode
	_ = g.AddEdge(PromptNode, ChatModelNode)      // PromptNode -> ChatModelNode
	_ = g.AddEdge(ChatModelNode, WriteOutputNode) // ChatModelNode -> WriteOutputNode
	_ = g.AddEdge(WriteOutputNode, compose.END)   // WriteOutputNode -> END

	// 编译 Graph
	r, err := g.Compile(ctx, compose.WithGraphName("InterviewAgentGraph"))
	if err != nil {
		return nil, err
	}
	return r, nil
}
