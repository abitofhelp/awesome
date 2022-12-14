package client_go

import (
	"context"
	"fmt"
	awesomev1 "github.com/abitofhelp/awesome/gen/go/awesome/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	kConnectionTimeOut = 1000 * time.Second
)

type PetReport struct {
	Report struct {
		AccessTier   string    `json:"accessTier"`
		Archived     bool      `json:"archived"`
		GeneratedUtc time.Time `json:"generatedUtc"`
		Pet          struct {
			BirthdayUtc time.Time `json:"birthdayUtc"`
			Name        string    `json:"name"`
		} `json:"pet"`
		Title string `json:"title"`
		URI   string `json:"uri"`
	} `json:"report"`
	Privacy string `json:"privacy"`
}

// AwesomeServiceClient implements the client_go.Persistence/gRPC AwesomeServiceClient interface.
type AwesomeServiceClient struct {
	client Persistence
}

// NewAwesomeServiceClient instantiates a new AwesomeServiceClient.
func NewAwesomeServiceClient(host string, port uint64) (*AwesomeServiceClient, error) {
	if len(host) == 0 {
		return nil, fmt.Errorf("the host cannot be empty or blank")
	}

	if port == 0 {
		return nil, fmt.Errorf("the port cannot be zero")
	}

	target := fmt.Sprintf("%s:%d", host, port)
	if conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		client := awesomev1.NewAwesomeServiceClient(conn)
		return &AwesomeServiceClient{client: client}, nil
	} else {
		return nil, fmt.Errorf("\nfailed to connect to the awesome service at '%s': %w", target, err)
	}
}

func (x *AwesomeServiceClient) FindReportByPetName(ctx context.Context, name string) (*PetReport, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("the name cannot be empty or blank")
	}

	ctx, cancel := context.WithTimeout(ctx, kConnectionTimeOut)
	defer cancel()

	if r, err := x.findReportByPetName(ctx, &awesomev1.FindReportByPetNameRequest{
		PetName: name,
	}); err == nil {

		report := &PetReport{
			Report: struct {
				AccessTier   string    `json:"accessTier"`
				Archived     bool      `json:"archived"`
				GeneratedUtc time.Time `json:"generatedUtc"`
				Pet          struct {
					BirthdayUtc time.Time `json:"birthdayUtc"`
					Name        string    `json:"name"`
				} `json:"pet"`
				Title string `json:"title"`
				URI   string `json:"uri"`
			}{
				AccessTier:   r.Report.AccessTier.String(),
				Archived:     r.Report.Archived,
				GeneratedUtc: r.Report.GeneratedUtc.AsTime(),
				Pet: struct {
					BirthdayUtc time.Time `json:"birthdayUtc"`
					Name        string    `json:"name"`
				}{
					BirthdayUtc: r.Report.Pet.BirthdayUtc.AsTime(),
					Name:        r.Report.Pet.Name,
				},
				Title: r.Report.Title,
				URI:   r.Report.Uri,
			},
			Privacy: r.Privacy.String(),
		}

		return report, nil
	} else {
		return nil, err
	}
}

func (x *AwesomeServiceClient) findReportByPetName(ctx context.Context, in *awesomev1.FindReportByPetNameRequest, opts ...grpc.CallOption) (*awesomev1.FindReportByPetNameResponse, error) {
	if r, err := x.client.FindReportByPetName(ctx, in); err == nil {
		return r, nil
	} else {
		return nil, fmt.Errorf("\nfailed to find the pet named '%s': %w", in.PetName, err)
	}
}
