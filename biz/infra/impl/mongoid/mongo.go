package mongoid

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/xh-polaris/synapse4b/biz/infra/contract/id"
)

func New(ctx context.Context) id.IDGenerator {
	return &objectIDGenerator{}
}

type objectIDGenerator struct{}

func (i *objectIDGenerator) GenID(_ context.Context) id.ID {
	return id.ID(bson.NewObjectID())
}

func (i *objectIDGenerator) GenMultiIDs(_ context.Context, counts int) (ids []id.ID) {
	for _ = range counts {
		ids = append(ids, id.ID(bson.NewObjectID()))
	}
	return
}
