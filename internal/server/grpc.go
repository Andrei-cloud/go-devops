package server

import (
	"context"
	"net"

	"github.com/andrei-cloud/go-devops/internal/hash"
	"github.com/andrei-cloud/go-devops/internal/model"
	pb "github.com/andrei-cloud/go-devops/internal/proto"
	"github.com/andrei-cloud/go-devops/internal/repo"
	"github.com/andrei-cloud/go-devops/internal/storage/filestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	repo   repo.Repository
	f      filestore.Filestore
	key    []byte
	subnet *net.IPNet
}

// NewMetricsServer - creates new server instance with all ingected dependencies for gRPC communications.
func NewMetricsServer(s server) *MetricsServer {
	return &MetricsServer{
		repo:   s.repo,
		f:      s.f,
		key:    s.key,
		subnet: s.subnet,
	}
}

// UpdateGauge - updates gauge metrics as gRPC request
func (s *MetricsServer) UpdateGauge(ctx context.Context, req *pb.UpdGaugeRequest) (*pb.UpdGaugeResponse, error) {
	var response pb.UpdGaugeResponse

	lm := model.Metric{
		ID:    req.Metric.Id,
		Delta: &req.Metric.Delta,
		Value: &req.Metric.Value,
		Hash:  req.Metric.Hash,
	}
	switch req.Metric.Mtype {
	case pb.Metric_COUNTER:
		lm.MType = "counter"
	case pb.Metric_GAUGE:
		lm.MType = "gauge"
	}

	valid, err := hash.Validate(lm, s.key)
	if err != nil {
		log.Debug().AnErr("Validate", err).Msg("UpdateBulkPost")
		return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`, req.Metric.Id)
	}

	if valid {

		if err := s.repo.UpdateGauge(ctx, req.Metric.Id, req.Metric.Value); err != nil {
			log.Error().AnErr("UpdateGauge", err).Msg("failed to update in repository")
			return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`, req.Metric.Id)
		}
	} else {
		return nil, status.Errorf(codes.FailedPrecondition, `invalid hash on metric: %s`, req.Metric.Id)
	}

	return &response, nil
}

// UpdateCounter - updates counter metrics as gRPC request
func (s *MetricsServer) UpdateCounter(ctx context.Context, req *pb.UpdCounterRequest) (*pb.UpdCounterResponse, error) {
	var response pb.UpdCounterResponse

	lm := model.Metric{
		ID:    req.Metric.Id,
		Delta: &req.Metric.Delta,
		Value: &req.Metric.Value,
		Hash:  req.Metric.Hash,
	}
	switch req.Metric.Mtype {
	case pb.Metric_COUNTER:
		lm.MType = "counter"
	case pb.Metric_GAUGE:
		lm.MType = "gauge"
	}

	valid, err := hash.Validate(lm, s.key)
	if err != nil {
		log.Debug().AnErr("Validate", err).Msg("UpdateBulkPost")
		return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`)
	}

	if valid {
		if err := s.repo.UpdateCounter(ctx, req.Metric.Id, req.Metric.Delta); err != nil {
			log.Error().AnErr("UpdateCounter", err).Msg("failed to update in repository")
			return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`, req.Metric.Id)
		}
	} else {
		return nil, status.Errorf(codes.FailedPrecondition, `invalid hash on metric: %s`, req.Metric.Id)
	}

	return &response, nil
}

// UpdateCounter - updates metrics in bulk as gRPC request.
func (s *MetricsServer) UpdateMetrics(ctx context.Context, req *pb.UpdMetricsRequest) (*pb.UpdMetricsResponse, error) {
	var response pb.UpdMetricsResponse

	for _, m := range req.Metrics {
		lm := model.Metric{
			ID:    m.Id,
			Delta: &m.Delta,
			Value: &m.Value,
			Hash:  m.Hash,
		}
		switch m.Mtype {
		case pb.Metric_COUNTER:
			lm.MType = "counter"
		case pb.Metric_GAUGE:
			lm.MType = "gauge"
		}

		valid, err := hash.Validate(lm, s.key)
		if err != nil {
			log.Error().AnErr("Validate", err).Msg("UpdateMetrics")
			return nil, status.Errorf(codes.FailedPrecondition, `Failed to update metric: %s`)
		}

		if valid {
			switch m.Mtype {
			case pb.Metric_GAUGE:
				if err := s.repo.UpdateGauge(ctx, m.Id, m.Value); err != nil {
					log.Error().AnErr("UpdateGauge", err).Msg("failed to update in repository")
					return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`, m.Id)
				}
			case pb.Metric_COUNTER:
				if err := s.repo.UpdateCounter(ctx, m.Id, m.Delta); err != nil {
					log.Error().AnErr("UpdateCounter", err).Msg("failed to update in repository")
					return nil, status.Errorf(codes.Internal, `Failed to update metric: %s`, m.Id)
				}
			default:
				return nil, status.Errorf(codes.Unimplemented, `unimplemented metric: %s`, m.Id)
			}
		} else {
			return nil, status.Errorf(codes.FailedPrecondition, `invalid hash on metric: %s`, m.Id)
		}
	}

	return &response, nil
}
