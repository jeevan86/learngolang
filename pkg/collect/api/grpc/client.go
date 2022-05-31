package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient(serverAddr string) pb.CollectClient {
	// https://zhuanlan.zhihu.com/p/530582571
	conn, err := grpc.DialContext(
		context.Background(),
		serverAddr,
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewCollectClient(conn)
}
