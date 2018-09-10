// client project main.go
package main

import (
	"flag"
	"log"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	endpoint = flag.String("endpoint", "127.0.0.1:8001", "xDS endpoint")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	//	ctx, _ := context.WithCancel(context.Background())
	ctx := context.Background()

	disclient := xdsapi.NewClusterDiscoveryServiceClient(conn)

	resp, err := disclient.FetchClusters(ctx, &xdsapi.DiscoveryRequest{
		TypeUrl: "type.googleapis.com/envoy.api.v2.Cluster",
		Node:    &envoy_api_v2_core1.Node{},
	})
	if err != nil {
		log.Fatalf("get Clusters resp error: %v", err)
		return
	}
	clusters := make([]*xdsapi.Cluster, 0)
	for _, res := range resp.Resources {
		cluster := xdsapi.Cluster{}
		cluster.Unmarshal(res.GetValue())
		clusters = append(clusters, &cluster)
	}

	log.Println("-----Clusters-------")
	for _, cluster := range clusters {
		log.Printf("%v", cluster)
	}

	lisclient := xdsapi.NewListenerDiscoveryServiceClient(conn)
	resp, err = lisclient.FetchListeners(ctx, &xdsapi.DiscoveryRequest{
		TypeUrl: "type.googleapis.com/envoy.api.v2.Listener",
		Node:    &envoy_api_v2_core1.Node{
		//			Id: c.ServiceNode,
		//Cluster: "default-cluster",
		},
	})
	if err != nil {
		log.Fatalf("get listeners resp error: %v", err)
		return
	}
	listeners := make([]*xdsapi.Listener, 0)
	for _, res := range resp.Resources {
		listener := xdsapi.Listener{}
		listener.Unmarshal(res.GetValue())
		listeners = append(listeners, &listener)
	}
	log.Println("-----listeners-------")
	for _, listener := range listeners {
		log.Printf("%v", listener)
	}

	rclient := xdsapi.NewRouteDiscoveryServiceClient(conn)
	resp, err = rclient.FetchRoutes(ctx, &xdsapi.DiscoveryRequest{
		TypeUrl: "type.googleapis.com/envoy.api.v2.RouteConfiguration",
		Node:    &envoy_api_v2_core1.Node{},
	})
	if err != nil {
		log.Fatalf("get routes resp error: %v", err)
		return
	}
	routes := make([]*xdsapi.RouteConfiguration, 0)
	for _, res := range resp.Resources {
		route := xdsapi.RouteConfiguration{}
		route.Unmarshal(res.GetValue())
		routes = append(routes, &route)
	}
	log.Println("-----routes-------")
	for _, route := range routes {
		log.Printf("%v", route)
	}

	eclient := xdsapi.NewEndpointDiscoveryServiceClient(conn)
	resp, err = eclient.FetchEndpoints(ctx, &xdsapi.DiscoveryRequest{
		TypeUrl: "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
		Node:    &envoy_api_v2_core1.Node{},
	})
	if err != nil {
		log.Fatalf("get Endpoint resp error: %v", err)
		return
	}

	lbAssignments := make([]*xdsapi.ClusterLoadAssignment, 0)
	for _, res := range resp.Resources {
		lbAssignment := xdsapi.ClusterLoadAssignment{}
		lbAssignment.Unmarshal(res.GetValue())
		lbAssignments = append(lbAssignments, &lbAssignment)
	}
	log.Println("-----Endpoints-------")
	for _, lbAssignment := range lbAssignments {
		log.Printf("%v", lbAssignment)
	}
}
