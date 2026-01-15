package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	shardingtable "github.com/envoyproxy/go-control-plane/examples/shardingtable/resource"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000

	ShardingTableTypeURL = "type.googleapis.com/shardingtable.v1.ShardingTable"
)

var (
	snapshotVersion int
)

type shardingTableServer struct {
	shardingtable.UnimplementedShardingDiscoveryServiceServer
	xdsServer xds.Server
}

func (s *shardingTableServer) StreamShardingTables(stream grpc.BidiStreamingServer[discoveryv3.DiscoveryRequest, discoveryv3.DiscoveryResponse]) error {
	return s.xdsServer.StreamHandler(stream, ShardingTableTypeURL)
}

func (s *shardingTableServer) FetchShardingTables(ctx context.Context, req *discoveryv3.DiscoveryRequest) (*discoveryv3.DiscoveryResponse, error) {
	return s.xdsServer.Fetch(ctx, req)
}

func (s *shardingTableServer) DeltaShardingTables(stream grpc.BidiStreamingServer[discoveryv3.DeltaDiscoveryRequest, discoveryv3.DeltaDiscoveryResponse]) error {
	return s.xdsServer.DeltaStreamHandler(stream, ShardingTableTypeURL)
}

func makeShardingTable(name, content string) *shardingtable.ShardingTable {
	return &shardingtable.ShardingTable{
		Name:    name,
		Content: content,
	}
}

func createSnapshot(version string) (*types.Snapshot, error) {
	tables := []types.SnapshotResource{
		{
			Name:     "table1",
			Resource: makeShardingTable("table1", fmt.Sprintf("Content for table1 - version %s - timestamp: %d", version, time.Now().Unix())),
		},
		{
			Name:     "table2",
			Resource: makeShardingTable("table2", fmt.Sprintf("Content for table2 - version %s - timestamp: %d", version, time.Now().Unix())),
		},
	}

	return types.NewSnapshot(version, map[string][]types.SnapshotResource{
		ShardingTableTypeURL: tables,
	})
}

func main() {
	ctx := context.Background()

	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	server := xds.NewServer(ctx, snapshotCache, nil)

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	shardingServer := &shardingTableServer{xdsServer: server}
	shardingtable.RegisterShardingDiscoveryServiceServer(grpcServer, shardingServer)

	lis, err := net.Listen("tcp", ":18000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting ShardingTable xDS server on port 18000")

	go func() {
		nodeID := "test-node-1"
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			snapshotVersion++
			versionStr := fmt.Sprintf("v%d", snapshotVersion)

			log.Printf("Creating snapshot version %s\n", versionStr)

			snapshot, err := createSnapshot(versionStr)
			if err != nil {
				log.Printf("Error creating snapshot: %v\n", err)
				continue
			}

			if err := snapshotCache.SetSnapshot(ctx, nodeID, snapshot); err != nil {
				log.Printf("Error setting snapshot: %v\n", err)
			} else {
				log.Printf("Snapshot %s set successfully\n", versionStr)
			}

			<-ticker.C
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
