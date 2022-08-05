package grpc

import (
	"context"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	PrepareTestServer()
	m.Run()
	StopTestServer()
}

var nodeIp = "192.168.3.153"
var serverAddr = "localhost:50051"

var c = NewClient(serverAddr)

func Test_Save(t *testing.T) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := c.Save(ctx, &pb.NetStaticsReq{GatherIp: nodeIp})
	if err != nil {
		log.Fatalf("Could not save: %v", err)
	}
	log.Printf("Save success: %s", r.GetMessage())
}

func Test_LocalIp(t *testing.T) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := c.LocalIp(ctx, &pb.LocalIpReq{NodeIp: nodeIp})
	if err != nil {
		log.Fatalf("Retrieve ip list failed: %v", err)
	}
	log.Printf("Retrieved ip list: %v", r.Data.GetIpList())
}
