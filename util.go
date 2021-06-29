package mgorm

import (
	"context"
	"fmt"
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
