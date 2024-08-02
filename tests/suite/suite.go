package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"sso/internal/config"

	ssov1 "github.com/Ranik23/proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T                  
	Cfg        *config.Config  
	AuthClient ssov1.AuthClient 
}

const (
	grpcHost = "localhost"
)


func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("localhost:4040")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.TimeOut)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}
