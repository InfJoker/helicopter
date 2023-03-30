package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "helicopter/generated/proto"
	"helicopter/internal/config"
	"helicopter/internal/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Grpc struct {
	address string
	server  *helicopterServer
}

type helicopterServer struct {
	pb.UnimplementedHelicopterServer
	storage core.Storage
}

var _ pb.HelicopterServer = (*helicopterServer)(nil)

func newServer(storage core.Storage) *helicopterServer {
	server := new(helicopterServer)
	server.storage = storage
	return server
}

func NewGrpc(cfg config.Config, storage core.Storage) *Grpc {
	host, port := "localhost", 1337
	if cfg.GrpcServer.Host != "" {
		host = cfg.GrpcServer.Host
	}
	if cfg.GrpcServer.Port != 0 {
		port = cfg.GrpcServer.Port
	}
	address := fmt.Sprintf("%s:%d", host, port)

	g := new(Grpc)
	g.server = newServer(storage)
	g.address = address
	return g
}

func (g *Grpc) Run() error {
	lis, err := net.Listen("tcp", g.address)
	if err != nil {
		log.Fatalf("Couldn't start listening on address \"%s\": %v", g.address, err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterHelicopterServer(grpcServer, g.server)
	reflection.Register(grpcServer)

	log.Print("Running rpc server at ", g.address)
	return grpcServer.Serve(lis)
}

func node2pbNode(node core.Node) *pb.Node {
	return &pb.Node{
		Lseq:    node.Lseq,
		Parent:  node.Parent,
		Content: node.Value,
	}
}

func nodes2pbNodes(nodes []core.Node) []*pb.Node {
	var pbNodes []*pb.Node
	for _, node := range nodes {
		pbNodes = append(pbNodes, node2pbNode(node))
	}
	return pbNodes
}

func (g *helicopterServer) GetNodes(ctx context.Context, req *pb.GetNodesRequest) (*pb.GetNodesResponse, error) {
	if req.Root == "0" {
		err := status.Error(codes.PermissionDenied, "Root must be non-zero")
		return nil, err
	}
	nodes, err := g.storage.GetSubTreeNodes(ctx, req.Root, req.Last)
	if err != nil {
		return nil, status.Error(codes.Internal, "Something broke 😿")
	}
	log.Printf("Got nodes %v\n", nodes)
	return &pb.GetNodesResponse{Nodes: nodes2pbNodes(nodes)}, nil
}

func (g *helicopterServer) AddNode(ctx context.Context, req *pb.AddNodeRequest) (*pb.AddNodeResponse, error) {
	node, err := g.storage.CreateNode(ctx, req.Parent, req.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, "Something broke 😿")
	}
	return &pb.AddNodeResponse{Node: node2pbNode(node)}, nil
}
