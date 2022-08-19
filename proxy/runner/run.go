package runner

import (
	"context"
	"fmt"
	"github.com/abitofhelp/awesome/config"
	awesomev1 "github.com/abitofhelp/awesome/gen/go/awesome/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	//	proxyv5 "github.com/ingios/project/v5/gen/go/project/v5"
	logger "github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

func Run(appcfg *config.AppConfig) error {

	ctx, cancel := context.WithTimeout(context.Background(), appcfg.Runtime.ConnectionTimeOut)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(
	//	runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
	//		// Attempted to get the Authorization header from the HTTP request.  If it exists,
	//		// we create an entry in the context's metadata so the downstream gRPC endpoints
	//		// can validate the access token.
	//		header := request.Header.Get("Authorization")
	//		// "Bearer aabcd..."
	//		md := metadata.Pairs("bearer", header)
	//		return md
	//	},
	//	)
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := awesomev1.RegisterAwesomeServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", appcfg.Runtime.Host, appcfg.Runtime.GrpcPort), opts)
	if err != nil {
		return err
	}

	con, err := net.Listen("tcp", fmt.Sprintf("%s:%d", appcfg.Runtime.Host, appcfg.Runtime.HttpPort))
	if err != nil {
		panic(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	logger.Printf(">>>> Starting HTTP proxy project service on %s...\n", con.Addr().String())
	return http.Serve(con, mux)
}
