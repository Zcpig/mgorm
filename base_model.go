package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type baseModel struct {
	sessCtx mongo.SessionContext
}

func (c *baseModel) Context() context.Context {
	if c.sessCtx != nil {
		return c.sessCtx
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func (c *baseModel) NewFilter(id string) bson.M {
	objId, _ := primitive.ObjectIDFromHex(id)
	return bson.M{"_id": objId}
}

func (c *baseModel) KeysMap(m bson.M) []string {
	var keys = make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
func (c *baseModel) SortMap(keys []string) bson.M {
	var m = bson.M{}
	for _, k := range keys {
		key := k
		val := 1
		if k[0] == '-' {
			key = k[1:]
			val = -1
		}
		m[key] = val
	}
	return m
}

func (c *baseModel) SetSessionContext(sessCtx mongo.SessionContext) {
	c.sessCtx = sessCtx
}
