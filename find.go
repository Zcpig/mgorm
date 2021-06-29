package mgorm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindBuilder struct {
	ctx    context.Context
	col    *mongo.Collection
	opt    *options.FindOptions
	filter interface{}
}

func (c *FindBuilder) Select(m bson.M) *FindBuilder {
	c.opt.SetProjection(m)
	return c
}

func (c *FindBuilder) Offset(num int64) *FindBuilder {
	c.opt.SetSkip(num)
	return c
}

func (c *FindBuilder) Limit(num int64) *FindBuilder {
	c.opt.SetLimit(num)
	return c
}

func (c *FindBuilder) Sort(sort interface{}) *FindBuilder {
	c.opt.SetSort(sort)
	return c
}

func (c *FindBuilder) Page(page, count int64) *FindBuilder {
	c.opt.SetLimit(count)
	c.opt.SetSkip((page - 1) * count)
	return c
}

func (c *FindBuilder) All(v interface{}) error {
	cursor, err := c.col.Find(c.ctx, c.filter, c.opt)
	if err != nil {
		return errorWrapper(err, c.col.Name())
	}
	return errorWrapper(cursor.All(c.ctx, v), executeStatement(c.col.Name(), c.filter, ""))
}

func (c *FindBuilder) One(v interface{}) error {
	var opt = options.FindOne()
	if c.opt.Projection != nil {
		opt.SetProjection(c.opt.Projection)
	}
	return errorWrapper(c.col.FindOne(c.ctx, c.filter, opt).Decode(v), executeStatement(c.col.Name(), c.filter, ""))
}
