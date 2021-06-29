package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateBuilder struct {
	ctx    context.Context
	opt    *options.UpdateOptions
	col    *mongo.Collection
	filter interface{}
	update interface{}
}

func (c *UpdateBuilder) SetUpsert(f bool) {
	c.opt.SetUpsert(f)
}

func (c *UpdateBuilder) One() (*mongo.UpdateResult, error) {
	result, err := c.col.UpdateOne(c.ctx, c.filter, c.update, c.opt)
	if err != nil {
		return nil, errorWrapper(err, executeStatement(c.col.Name(), c.filter, c.update))
	}
	return result, nil
}

func (c *UpdateBuilder) All() (*mongo.UpdateResult, error) {
	result, err := c.col.UpdateMany(c.ctx, c.filter, c.update, c.opt)
	if err != nil {
		return nil, errorWrapper(err, executeStatement(c.col.Name(), c.filter, c.update))
	}
	return result, nil
}
