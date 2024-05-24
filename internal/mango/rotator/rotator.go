package rotator

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

type ApiGateway struct {
	client    *apigatewayv2.Client
	site      string
	regions   []string
	endpoints []string
	accessKey string
	secretKey string
	mutex     sync.Mutex
}

func NewApiGateway(site string, regions []string, accessKey, secretKey string) *ApiGateway {
	client := apigatewayv2.New(apigatewayv2.Options{})

	return &ApiGateway{
		client:    client,
		site:      site,
		regions:   regions,
		endpoints: []string{},
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// Fetches endpoints from different regions (replace with actual API Gateway logic)
func (ag *ApiGateway) fetchEndpoints() error {
	ag.mutex.Lock()
	defer ag.mutex.Unlock()

	// ... (Simulate fetching endpoints from each region using AWS SDK)

	return nil
}

// Starts the API Gateway by fetching endpoints (modify for your use case)
func (ag *ApiGateway) Start() error {
	err := ag.fetchEndpoints()
	if err != nil {
		return err
	}

	fmt.Printf("Using %d endpoints for site '%s'\n", len(ag.endpoints), ag.site)
	return nil
}
