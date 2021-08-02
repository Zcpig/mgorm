package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type BaseModel struct {
	sessCtx mongo.SessionContext
}

func (c *BaseModel) Context() context.Context {
	if c.sessCtx != nil {
		return c.sessCtx
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func (c *BaseModel) NewFilter(id string) bson.M {
	objId, _ := primitive.ObjectIDFromHex(id)
	return bson.M{"_id": objId}
}

func (c *BaseModel) KeysMap(m bson.M) []string {
	var keys = make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
func (c *BaseModel) SortMap(fields []string) interface{} {
	var order bson.D
	for _, field := range fields {
		n := 1
		var kind string
		if field != "" {
			if field[0] == '$' {
				if c := strings.Index(field, ":"); c > 1 && c < len(field)-1 {
					kind = field[1:c]
					field = field[c+1:]
				}
			}
			switch field[0] {
			case '+':
				field = field[1:]
			case '-':
				n = -1
				field = field[1:]
			}
		}
		if field == "" {
			panic("Sort: empty field name")
		}
		if kind == "textScore" {
			order = append(order, bson.E{Key: field, Value: bson.M{"$meta": kind}})
		} else {
			order = append(order, bson.E{Key: field, Value: n})
		}
	}
	return order
}

func (c *BaseModel) SetSessionContext(sessCtx mongo.SessionContext) {
	c.sessCtx = sessCtx
}
