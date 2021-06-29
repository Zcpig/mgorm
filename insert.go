package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertBuilder struct {
	ctx context.Context
	col *mongo.Collection
}

func (c *InsertBuilder) One(document interface{}) (primitive.ObjectID, error) {
	result, err := c.col.InsertOne(c.ctx, document)
	if err != nil {
		return primitive.ObjectID{}, errorWrapper(err, executeStatement(c.col.Name(), "", document))
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (c *InsertBuilder) All(documents []interface{}) ([]primitive.ObjectID, error) {
	results, err := c.col.InsertMany(c.ctx, documents)
	if err != nil {
		return nil, errorWrapper(err, executeStatement(c.col.Name(), "", documents))
	}

	var ids = make([]primitive.ObjectID, 0)
	for _, item := range results.InsertedIDs {
		ids = append(ids, item.(primitive.ObjectID))
	}
	return ids, nil
}
