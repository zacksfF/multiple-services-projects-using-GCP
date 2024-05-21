package tracer

import (
	"github.com/zacksfF/Build-A-Go-Apps-That-Scale-On-Google-Cloud/Observability/logging"
	"go.opentelemetry.io/otel/sdk/trace"
)

type endpointExcluder struct {
	log         *logging.Logger
	endpoints   map[string]struct{}
	probability float64
}

func newEndpointExcluder(log *logging.Logger, endpoints map[string]struct{}, probability float64) endpointExcluder {
	return endpointExcluder{
		log:         log,
		endpoints:   endpoints,
		probability: probability,
	}
}

// ShouldSample implements the sampler interface. It prevents the specified
// endpoints from being added to the trace.
func (ee endpointExcluder) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	for i := range parameters.Attributes {
		if parameters.Attributes[i].Key == "http.target" {
			if _, exists := ee.endpoints[parameters.Attributes[i].Value.AsString()]; exists {
				return trace.SamplingResult{Decision: trace.Drop}
			}
		}
	}

	return trace.TraceIDRatioBased(ee.probability).ShouldSample(parameters)
}

// Description implements the sampler interface.
func (endpointExcluder) Description() string {
	return "customSampler"
}
