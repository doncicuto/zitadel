package schemauser

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "user"
	AggregateVersion = "v3"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, resourceOwner string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: resourceOwner,
		},
	}
}
