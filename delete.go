package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeleteBuilder struct {
	ctx    context.Context
	filter interface{}
	opt    *options.DeleteOptions
	col    *mongo.Collection
}

func (c *DeleteBuilder) One() (*mongo.DeleteResult, error) {
	result, err := c.col.DeleteOne(c.ctx, c.filter, c.opt)
	if err != nil {
		return nil, errorWrapper(err, executeStatement(c.col.Name(), c.filter, ""))
	}
	return result, nil
}

func (c *DeleteBuilder) All() (*mongo.DeleteResult, error) {
	result, err := c.col.DeleteMany(c.ctx, c.filter, c.opt)
	if err != nil {
		return nil, errorWrapper(err, executeStatement(c.col.Name(), c.filter, ""))
	}
	return result, nil
}
