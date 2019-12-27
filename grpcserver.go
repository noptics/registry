package main

import (
	"context"
	"net"

	"github.com/noptics/golog"
	"github.com/noptics/registry/data"
	"github.com/noptics/registry/registrygrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	l       golog.Logger
	db      data.Store
	grpcs   *grpc.Server
	errChan chan error
}

func NewGRPCServer(db data.Store, port string, errChan chan error, l golog.Logger) (*GRPCServer, error) {
	gs := &GRPCServer{
		db:      db,
		l:       l,
		errChan: errChan,
		grpcs:   grpc.NewServer(),
	}

	registrygrpc.RegisterProtoRegistryServer(gs.grpcs, gs)

	l.Infow("starting grpc service", "port", port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	go func() {
		if err = gs.grpcs.Serve(lis); err != nil {
			gs.errChan <- err
		}
	}()

	return gs, nil
}

func (s *GRPCServer) Stop() {
	s.grpcs.GracefulStop()
}

func (s *GRPCServer) GetFiles(ctx context.Context, in *registrygrpc.GetFilesRequest) (*registrygrpc.GetFilesReply, error) {
	_, f, err := s.db.GetChannelData(in.Cluster, in.Channel)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.GetFilesReply{Files: f}, nil
}

func (s *GRPCServer) SaveFiles(ctx context.Context, in *registrygrpc.SaveFilesRequest) (*registrygrpc.SaveFilesReply, error) {
	err := s.db.SaveFiles(in.Cluster, in.Channel, in.Files)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.SaveFilesReply{}, nil
}

func (s *GRPCServer) SetMessage(ctx context.Context, in *registrygrpc.SetMessageRequest) (*registrygrpc.SetMessageReply, error) {
	err := s.db.SetChannelMessage(in.Cluster, in.Channel, in.Name)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.SetMessageReply{}, nil
}

func (s *GRPCServer) GetMessage(ctx context.Context, in *registrygrpc.GetMessageRequest) (*registrygrpc.GetMessageReply, error) {
	m, _, err := s.db.GetChannelData(in.Cluster, in.Channel)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.GetMessageReply{Name: m}, nil
}

func (s *GRPCServer) GetChannelData(ctx context.Context, in *registrygrpc.GetChannelDataRequest) (*registrygrpc.GetChannelDataReply, error) {
	m, files, err := s.db.GetChannelData(in.Cluster, in.Channel)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.GetChannelDataReply{Cluster: in.Cluster, Channel: in.Channel, Files: files, Message: m}, nil
}

func (s *GRPCServer) SetChannelData(ctx context.Context, in *registrygrpc.SetChannelDataRequest) (*registrygrpc.SetChannelDataReply, error) {
	err := s.db.SaveChannelData(in.Cluster, in.Channel, in.Message, in.Files)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.SetChannelDataReply{}, nil
}

func (s *GRPCServer) SaveCluster(ctx context.Context, in *registrygrpc.SaveClusterRequest) (*registrygrpc.SaveClusterReply, error) {
	id, err := s.db.SaveCluster(in.Cluster)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.SaveClusterReply{Id: id}, nil
}

func (s *GRPCServer) GetCluster(ctx context.Context, in *registrygrpc.GetClusterRequest) (*registrygrpc.GetClusterReply, error) {
	cluster, err := s.db.GetCluster(in.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.GetClusterReply{Cluster: cluster}, nil
}

func (s *GRPCServer) GetClusters(ctx context.Context, in *registrygrpc.GetClustersRequest) (*registrygrpc.GetClustersReply, error) {
	clusters, err := s.db.GetClusters()

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &registrygrpc.GetClustersReply{Clusters: clusters}, nil
}
