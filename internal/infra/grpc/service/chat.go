package service

import (
	"github.com/lacerda.jcarlos/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/lacerda.jcarlos/fclx/chatservice/internal/usecase/chatcompletionstream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigInputDTO
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewChatService(chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase,
	chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO,
	streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}
}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompletionstream.ChatCompletionConfigInputDTO{
		Model:                c.ChatConfigStream.Model,
		ModelMaxTokens:       c.ChatConfigStream.ModelMaxTokens,
		Temperature:          c.ChatConfigStream.Temperature,
		TopP:                 c.ChatConfigStream.TopP,
		N:                    c.ChatConfigStream.N,
		Stop:                 c.ChatConfigStream.Stop,
		MaxTokens:            c.ChatConfigStream.MaxTokens,
		InitialSystemMessage: c.ChatConfigStream.InitialSystemMessage,
	}

	input := chatcompletionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()
	errChan := make(chan error)

	go func() {
		defer close(c.StreamChannel)
		for msg := range c.StreamChannel {
			err := stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	_, err := c.ChatCompletionStreamUseCase.Execute(ctx, input)
	if err != nil {
		return err
	}

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
