package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	shardingtable "github.com/envoyproxy/go-control-plane/examples/shardingtable/resource"
)

const (
	serverAddress = "localhost:18000"
	nodeID        = "test-node-1"
	typeURL       = "type.googleapis.com/shardingtable.v1.ShardingTable"
)

func main() {
	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := shardingtable.NewShardingDiscoveryServiceClient(conn)

	ctx := context.Background()

	log.Printf("Starting xDS delta stream to %s for node %s\n", serverAddress, nodeID)

	stream, err := client.DeltaShardingTables(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	initialReq := &discoveryv3.DeltaDiscoveryRequest{
		Node: &corev3.Node{
			Id:      nodeID,
			Cluster: "test-cluster",
		},
		TypeUrl:                 typeURL,
		ResourceNamesSubscribe:  []string{"table1", "table2"},
		InitialResourceVersions: make(map[string]string),
	}

	log.Println("Sending initial discovery request...")
	if err := stream.Send(initialReq); err != nil {
		log.Fatalf("Failed to send initial request: %v", err)
	}

	resourceVersions := make(map[string]string)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Stream closed by server")
				return
			}
			if err != nil {
				log.Printf("Error receiving response: %v", err)
				return
			}

			log.Printf("\n=== Received DeltaDiscoveryResponse ===")
			log.Printf("System Version: %s", resp.SystemVersionInfo)
			log.Printf("Nonce: %s", resp.Nonce)
			log.Printf("Type URL: %s", resp.TypeUrl)
			log.Printf("Resources count: %d", len(resp.Resources))
			log.Printf("Removed resources count: %d", len(resp.RemovedResources))

			for i, resource := range resp.Resources {
				table := &shardingtable.ShardingTable{}
				if err := resource.Resource.UnmarshalTo(table); err != nil {
					log.Printf("Failed to unmarshal resource %d: %v", i, err)
					continue
				}

				log.Printf("\nShardingTable %d:", i+1)
				log.Printf("  Name: %s", table.Name)
				log.Printf("  Content: %s", table.Content)
				log.Printf("  Resource Version: %s", resource.Version)

				resourceVersions[resource.Name] = resource.Version
			}

			for _, removedName := range resp.RemovedResources {
				log.Printf("\nResource removed: %s", removedName)
				delete(resourceVersions, removedName)
			}

			log.Println("==================================\n")

			ackReq := &discoveryv3.DeltaDiscoveryRequest{
				Node: &corev3.Node{
					Id:      nodeID,
					Cluster: "test-cluster",
				},
				TypeUrl:                 typeURL,
				ResponseNonce:           resp.Nonce,
				InitialResourceVersions: resourceVersions,
			}

			if err := stream.Send(ackReq); err != nil {
				log.Printf("Failed to send ACK: %v", err)
				return
			}
			log.Printf("Sent ACK for nonce %s with %d tracked resources\n", resp.Nonce, len(resourceVersions))
		}
	}()

	log.Println("Client is running. Press Ctrl+C to exit.")
	select {}
}
