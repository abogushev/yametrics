package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	pb "yametrics/internal/protocol/proto"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
}

func RunMetricsServer(logger *zap.SugaredLogger, ctx context.Context) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Error(err)
		return
	}
	// создаём gRPC-сервер без зарегистрированной службы
	server := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(server, &MetricsServer{})

	logger.Info("Сервер gRPC начал работу")
	// получаем запрос gRPC
	go func() {
		if err := server.Serve(listen); err != nil {
			logger.Error(err)
		}
	}()
	<-ctx.Done()
	server.GracefulStop()
	logger.Info("server stopped")
}

func (s *MetricsServer) SaveMetrics(stream pb.Metrics_SaveMetricsServer) error {
	for {
		metric, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}
		fmt.Println(metric.Id)
	}
}
