package runner

import (
	"context"
	"fmt"
	"github.com/abitofhelp/awesome/config"
	awesomev1 "github.com/abitofhelp/awesome/gen/go/awesome/v1"
	enumsv1 "github.com/abitofhelp/awesome/gen/go/awesome/v1/enums"
	messagesv1 "github.com/abitofhelp/awesome/gen/go/awesome/v1/messages"
	logger "github.com/labstack/gommon/log"
	//enumsv1 "go.buf.build/grpc/go/abitofhelp/abcdapis/enums/v1"
	//messagesv1 "go.buf.build/grpc/go/abitofhelp/abcdapis/messages/v1"
	//awesomev1 "go.buf.build/grpc/go/abitofhelp/abcdapis/services/awesome/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

func Run(appcfg *config.AppConfig) error {

	server := grpc.NewServer()

	svc, err := NewService()
	if err != nil {
		panic(err)
	}

	awesomev1.RegisterAwesomeServiceServer(server, svc)
	reflection.Register(server)

	con, err := net.Listen("tcp", fmt.Sprintf("%s:%d", appcfg.Runtime.Host, appcfg.Runtime.GrpcPort))
	if err != nil {
		panic(err)
	}

	logger.Printf(">>>> Starting gRPC project service on %s...\n", con.Addr().String())
	return server.Serve(con)
}

type Server struct {
	awesomev1.AwesomeServiceServer
}

func (s Server) FindReportByPetName(
	ctx context.Context,
	in *awesomev1.FindReportByPetNameRequest) (*awesomev1.FindReportByPetNameResponse, error) {
	return s.findReportByPetName(ctx, in)
}

func (s Server) findReportByPetName(
	ctx context.Context,
	in *awesomev1.FindReportByPetNameRequest,
	opts ...grpc.CallOption) (*awesomev1.FindReportByPetNameResponse, error) {

	return &awesomev1.FindReportByPetNameResponse{
		Report: &messagesv1.Report{
			AccessTier:   enumsv1.AccessTier_ACCESS_TIER_COOL,
			Archived:     true,
			GeneratedUtc: timestamppb.New(time.Now().UTC()),
			Pet: &messagesv1.Pet{
				BirthdayUtc: timestamppb.New(time.Now().UTC()),
				Name:        "Lassie",
			},
			Title: "Lassie's Report",
			Uri:   "http://reports.com/lassie",
		},
		Privacy: awesomev1.Privacy_PRIVACY_NONE,
	}, nil
}

func NewService() (*Server, error) {
	return &Server{}, nil
}
