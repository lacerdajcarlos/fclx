package server

import (
	"log"
	"net"

	"github.com/lacerda.jcarlos/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/lacerda.jcarlos/fclx/chatservice/internal/infra/grpc/service"
	"github.com/lacerda.jcarlos/fclx/chatservice/internal/usecase/chatcompletionstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigInputDTO
	ChatService                 *service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewGRPCServer(
	chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase,
	chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO,
	port string,
	authToken string,
	streamChannel chan chatcompletionstream.ChatCompletionOutputDTO,
) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)
	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               streamChannel,
		ChatService:                 chatService,
	}
}

// AuthInterceptor para streams
func (g *GRPCServer) AuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	log.Printf("Received token: %s", token[0])    // Log do token recebido
	log.Printf("Expected token: %s", g.AuthToken) // Log do token esperado

	if token[0] != g.AuthToken {
		return status.Error(codes.Unauthenticated, "authorization token is invalid")
	}

	log.Printf("Client authenticated successfully")
	return handler(srv, ss)
}

// Método para iniciar o servidor gRPC
func (g *GRPCServer) Start() {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(g.AuthInterceptor),
	}

	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", g.Port, err)
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServiceServer(grpcServer, g.ChatService)

	// Habilita a reflexão gRPC
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on port %s", g.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
