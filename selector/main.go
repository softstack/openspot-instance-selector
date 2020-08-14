package main

import (
	"context"
	"fmt"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/bytequantity"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SelectorResponse struct {
	Instances []string `json:"instances"`
}

type SelectorEvent struct {
	MinVcpus int `json:"min_cpus, omitempty"`
	MaxVcpus int `json:"max_cpus, omitempty"`
	MinMemory uint64 `json:"min_memory, omitempty"`
	MaxMemory uint64 `json:"max_memory, omitempty"`
	MinNetwork int `json:"min_network, omitempty"`
	MaxNetwork int `json:"max_network, omitempty"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event SelectorEvent) (*SelectorResponse, error) {

	sess, err := session.NewSession(nil)

	if err != nil {
		return nil, err
	}

	instanceSelector := selector.New(sess)

	vcpusRange := selector.IntRangeFilter{
		LowerBound: event.MinVcpus,
		UpperBound: event.MaxVcpus,
	}

	memoryRange := selector.ByteQuantityRangeFilter{
		LowerBound: bytequantity.FromMiB(event.MinMemory),
		UpperBound: bytequantity.FromMiB(event.MaxMemory),
	}

	networkRange := selector.IntRangeFilter{
		LowerBound: event.MinNetwork,
		UpperBound: event.MaxNetwork,
	}

	filters := selector.Filters{
		VCpusRange: &vcpusRange,
		MemoryRange: &memoryRange,
		NetworkPerformance: &networkRange,
	}

	instanceTypesSlice, err := instanceSelector.Filter(filters)

	if err != nil {
		fmt.Printf("Oh no, there was an error :( %v", err)
		return nil, err
	}

	resp := &SelectorResponse{Instances: instanceTypesSlice}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
