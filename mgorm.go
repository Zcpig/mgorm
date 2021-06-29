package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errorWrapper = func(err error, colName string) error {
	return err
}

func SetErrorWrapper(fn func(err error, colName string) error) {
	errorWrapper = fn
}

type MgORM struct {
	col *mongo.Collection
}

type Option struct {
	Collection string
	WrapError  func(err error) error
}

func NewORM(collection *mongo.Collection) *MgORM {
	return &MgORM{
		col: collection,
	}
}

func (c *MgORM) Collection() *mongo.Collection {
	return c.col
}

func (c *MgORM) Find(ctx context.Context, filter interface{}) *FindBuilder {
	if filter == nil {
		filter = bson.M{}
	}
	return &FindBuilder{
		ctx:    ctx,
		col:    c.col,
		opt:    options.Find(),
		filter: filter,
	}
}

func (c *MgORM) Count(ctx context.Context, filter interface{}) (int64, error) {
	count, err := c.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errorWrapper(err, c.col.Name())
	}
	return count, nil
}

func (c *MgORM) Exist(ctx context.Context, filter interface{}) bool {
	count, err := c.col.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}
	return count > 0
}

func (c *MgORM) Update(ctx context.Context, filter interface{}, update interface{}) *UpdateBuilder {
	return &UpdateBuilder{
		ctx:    ctx,
		opt:    options.Update(),
		col:    c.col,
		filter: filter,
		update: update,
	}
}

func (c *MgORM) Delete(ctx context.Context, filter interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		ctx:    ctx,
		opt:    options.Delete(),
		col:    c.col,
		filter: filter,
	}
}

func (c *MgORM) Insert(ctx context.Context) *InsertBuilder {
	return &InsertBuilder{ctx: ctx, col: c.col}
}

func (c *MgORM) NewTransaction(ctx context.Context, callback func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := c.col.Database().Client().StartSession()
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	return result, nil
}

func (c *MgORM) CreateIndex(keys []string, name string, unique bool) error {
	opt := options.Index().SetUnique(unique)
	if name != "" {
		opt.SetName(name)
	}

	var bsonKeys = bson.D{}
	for _, key := range keys {
		bsonKeys = append(bsonKeys, bson.E{Key: key, Value: 1})
	}

	_, err := c.col.Indexes().CreateOne(Context(), mongo.IndexModel{
		Keys:    bsonKeys,
		Options: opt.SetBackground(true),
	})
	return errorWrapper(err, c.col.Name())
}

func (c *MgORM) FindById(ctx context.Context, id string, result interface{}) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errorWrapper(err, c.col.Name())
	}
	var filter = bson.M{"_id": objId}
	return errorWrapper(c.col.FindOne(ctx, filter).Decode(result), c.col.Name())
}

func (c *MgORM) UpdateById(ctx context.Context, id string, update interface{}) (*mongo.UpdateResult, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	var filter = bson.M{"_id": objId}
	result, err := c.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	return result, nil
}

func (c *MgORM) DeleteById(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	var filter = bson.M{"_id": objId}
	result, err := c.col.DeleteOne(ctx, filter)
	if err != nil {
		return nil, errorWrapper(err, c.col.Name())
	}
	return result, nil
}

func (c *MgORM) Aggregate(ctx context.Context) *AggregateBuilder {
	return &AggregateBuilder{
		ctx:      ctx,
		col:      c.col,
		pipeline: bson.A{},
	}
}
