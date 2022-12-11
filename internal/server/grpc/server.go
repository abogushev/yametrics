package grpc

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	pb "yametrics/internal/protocol/proto"
	"yametrics/internal/server/models"
	"yametrics/internal/server/storage"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	logger         *zap.SugaredLogger
	metricsStorage storage.MetricsStorage
}

func RunMetricsServer(logger *zap.SugaredLogger, ctx context.Context, metricsStorage storage.MetricsStorage) {
	creds, err := credentials.NewServerTLSFromFile("cert/service.pem", "cert/service.key")
	if err != nil {
		logger.Fatalf("Failed to setup TLS: %v", err)
	}
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		logger.Error(err)
		return
	}
	// создаём gRPC-сервер без зарегистрированной службы
	server := grpc.NewServer(grpc.Creds(creds))
	// регистрируем сервис
	pb.RegisterMetricsServer(server, &MetricsServer{logger: logger, metricsStorage: metricsStorage})

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
	mtrcs := make([]models.Metrics, 0)

	for {
		metric, err := stream.Recv()

		if err == io.EOF {
			err = s.metricsStorage.Updates(mtrcs)
			if err != nil {
				return err
			}
			s.logger.Info("metrics saved successful")
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}
		t := ""
		switch metric.Type {
		case 0:
			t = "counter"
		case 1:
			t = "gauge"
		}
		mtrcs = append(mtrcs, models.Metrics{
			ID:    metric.Id,
			MType: t,
			Delta: metric.Delta,
			Value: metric.Value,
		})
	}
}
