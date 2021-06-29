package mgorm

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func WithTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}

// default context
func Context() context.Context {
	return WithTimeout(15 * time.Second)
}

func executeStatement(colName string, filter, update interface{}) string {
	return fmt.Sprintf("colName:%s - filter:%v - update :%v", colName, filter, update)
}

func IsObjectIdHex(s string) bool {
	_, err := primitive.ObjectIDFromHex(s)
	return err == nil
}

func ObjectIdHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.ObjectID{}
	}
	return id
}
