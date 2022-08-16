package agent

import (
	"context"
	"fmt"

	"github.com/andrei-cloud/go-devops/internal/hash"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/andrei-cloud/go-devops/internal/proto"
)

// ReportCounterPost - reports counter metric to the sever.
func (a *agent) ReportCounterGRPC(ctx context.Context, m map[string]int64) {
	ipAddr := getLocalIP()
	log.Debug().Msgf("Real IP: %v", ipAddr)

	md := metadata.New(map[string]string{"X-Real-IP": ipAddr})
	lctx := metadata.NewOutgoingContext(ctx, md)

	for k, v := range m {
		metric := pb.Metric{}
		metric.Id = k
		metric.Mtype = pb.Metric_COUNTER
		metric.Delta = v

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:counter:%d", metric.Id, metric.Delta), a.key)
		}

		_, err := a.gclient.UpdateCounter(lctx, &pb.UpdCounterRequest{
			Metric: &metric,
		})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.Internal {
					log.Warn().Msgf(`INTERNAL SERVER ERROR: %s`, e.Message())
				} else {
					log.Error().Msgf("Error: %s - Message: %s", e.Code(), e.Message())
				}
			} else {
				log.Error().Msgf("Unable to parse error %v", err)
			}
		}
	}
}

func (a *agent) ReportGaugeGRPC(ctx context.Context, m map[string]float64) {
	ipAddr := getLocalIP()
	log.Debug().Msgf("Real IP: %v", ipAddr)

	md := metadata.New(map[string]string{"X-Real-IP": ipAddr})
	lctx := metadata.NewOutgoingContext(ctx, md)

	for k, v := range m {
		metric := pb.Metric{}
		metric.Id = k
		metric.Mtype = pb.Metric_GAUGE
		metric.Value = v

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:gauge:%f", metric.Id, metric.Value), a.key)
		}

		_, err := a.gclient.UpdateGauge(lctx, &pb.UpdGaugeRequest{
			Metric: &metric,
		})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.Internal {
					log.Warn().Msgf(`INTERNAL SERVER ERROR: %s`, e.Message())
				} else {
					log.Error().Msgf("Error: %s - Message: %s", e.Code(), e.Message())
				}
			} else {
				log.Error().Msgf("Unable to parse error %v", err)
			}
		}
	}
}

// ReportBulkPost - reports metrics in bulk to the sever.
func (a *agent) ReportBulkGRPC(ctx context.Context, c map[string]int64, g map[string]float64) {
	var req pb.UpdMetricsRequest

	ipAddr := getLocalIP()
	log.Debug().Msgf("Real IP: %v", ipAddr)

	md := metadata.New(map[string]string{"X-Real-IP": ipAddr})
	lctx := metadata.NewOutgoingContext(ctx, md)

	for k, v := range g {
		metric := pb.Metric{}

		locV := v
		metric.Id = k
		metric.Mtype = pb.Metric_GAUGE
		metric.Value = locV

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:gauge:%f", metric.Id, metric.Value), a.key)
		}
		req.Metrics = append(req.Metrics, &metric)
	}

	for k, v := range c {
		metric := pb.Metric{}

		locC := v
		metric.Id = k
		metric.Mtype = pb.Metric_COUNTER
		metric.Delta = locC

		if len(a.key) != 0 {
			metric.Hash = hash.Create(fmt.Sprintf("%s:counter:%d", metric.Id, metric.Delta), a.key)
		}
		req.Metrics = append(req.Metrics, &metric)
	}

	if len(req.Metrics) > 0 {
		_, err := a.gclient.UpdateMetrics(lctx, &req)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.Internal {
					log.Warn().Msgf(`INTERNAL SERVER ERROR: %s`, e.Message())
				} else {
					log.Error().Msgf("Error: %s - Message: %s", e.Code(), e.Message())
				}
			} else {
				log.Error().Msgf("Unable to parse error %v", err)
			}
		}
	}
}
