package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateBuilder struct {
	ctx      context.Context
	col      *mongo.Collection
	pipeline bson.A
}

func (c *AggregateBuilder) Match(match interface{}) *AggregateBuilder {
	var p = bson.M{
		"$match": match,
	}
	c.pipeline = append(c.pipeline, p)
	return c
}

func (c *AggregateBuilder) Project(project interface{}) *AggregateBuilder {
	var p = bson.M{
		"$project": project,
	}
	c.pipeline = append(c.pipeline, p)
	return c
}

func (c *AggregateBuilder) Limit(limit int64) *AggregateBuilder {
	var p = bson.M{
		"$limit": limit,
	}
	c.pipeline = append(c.pipeline, p)
	return c
}

func (c *AggregateBuilder) Offset(offset int64) *AggregateBuilder {
	var p = bson.M{
		"$skip": offset,
	}
	c.pipeline = append(c.pipeline, p)
	return c
}

func (c *AggregateBuilder) Group(key interface{}, m bson.M) *AggregateBuilder {
	m["_id"] = nil
	if id, ok := key.(string); ok {
		m["_id"] = "$" + id
	}

	var m1 = bson.M{
		"$group": m,
	}
	c.pipeline = append(c.pipeline, m1)
	return c
}

func (c *AggregateBuilder) Sort(sort interface{}) *AggregateBuilder {
	var p = bson.M{
		"$sort": sort,
	}
	c.pipeline = append(c.pipeline, p)
	return c
}

func (c *AggregateBuilder) One(result interface{}) error {
	cursor, err := c.col.Aggregate(c.ctx, c.pipeline)
	if err != nil {
		return errorWrapper(err, c.col.Name())
	}
	return errorWrapper(cursor.Decode(result), executeStatement(c.col.Name(), c.pipeline, ""))
}

func (c *AggregateBuilder) All(results interface{}) error {
	cursor, err := c.col.Aggregate(c.ctx, c.pipeline)
	if err != nil {
		return errorWrapper(err, c.col.Name())
	}
	return errorWrapper(cursor.All(c.ctx, results), executeStatement(c.col.Name(), c.pipeline, ""))
}
