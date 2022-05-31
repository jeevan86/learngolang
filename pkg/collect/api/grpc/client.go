package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient(serverAddr string) pb.CollectClient {
	//flag.Parse()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewCollectClient(conn)
}
