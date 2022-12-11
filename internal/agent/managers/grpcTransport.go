package managers

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"yametrics/internal/agent/models/storage"
	pb "yametrics/internal/protocol/proto"
)

type GRPCTransportManager struct {
	logger *zap.SugaredLogger
}

func NewGRPCTransportManager(logger *zap.SugaredLogger) *GRPCTransportManager {
	return &GRPCTransportManager{logger: logger}
}

func (t *GRPCTransportManager) Send(ctx context.Context, metrics *storage.Metrics) error {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.logger.Error(err)
	}
	defer conn.Close()
	c := pb.NewMetricsClient(conn)
	stream, err := c.SaveMetrics(ctx)
	if err != nil {
		t.logger.Error(err)
	}

	metricsForSend := make([]*pb.Metric, 0)

	metrics.OperateOverMetricMaps(
		func(s string, f float64) {
			metricsForSend = append(metricsForSend,
				&pb.Metric{
					Id:    s,
					Type:  pb.MetricTypes_GAUGE,
					Value: &f,
				})
		},
		func(s string, i int64) {
			metricsForSend = append(metricsForSend,
				&pb.Metric{
					Id:    s,
					Type:  pb.MetricTypes_COUNTER,
					Delta: &i,
				})
		})

	for i := 0; i < len(metricsForSend); i++ {
		err := stream.Send(metricsForSend[i])
		if err != nil {
			return err
		}
	}
	_, err = stream.CloseAndRecv()
	return err
}
