package dependences

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
)

type MongoDBClient struct {
	Client *mongo.Client
}

func (m *MongoDBClient) Disconnect(ctx context.Context) error {
	if m.Client != nil {
		err := m.Client.Disconnect(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoDBClient) InsertOne(ctx context.Context, database string, collection string, document interface{}) (map[string]interface{}, error) {
	result, err := m.Client.Database(database).Collection(collection).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	docMap := document.(map[string]interface{})
	id := result.InsertedID
	if oid, ok := id.(primitive.ObjectID); ok {
		docMap["_id"] = oid.Hex()
	} else {
		return nil, fmt.Errorf("inserted ID is not an ObjectID")
	}
	return docMap, nil
}

func (m *MongoDBClient) InsertMany(ctx context.Context, database, collection string, documents []interface{}) ([]map[string]interface{}, error) {
	result, err := m.Client.Database(database).Collection(collection).InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	var docs []map[string]interface{}
	for i, id := range result.InsertedIDs {
		docMap := documents[i].(map[string]interface{})
		if oid, ok := id.(primitive.ObjectID); ok {
			docMap["_id"] = oid.Hex()
		} else {
			return nil, fmt.Errorf("inserted ID is not an ObjectID")
		}
		docs = append(docs, docMap)
	}
	return docs, nil
}

func (m *MongoDBClient) FindOne(ctx context.Context, database, collection string, filter interface{}) (bson.M, error) {
	var result bson.M
	err := m.Client.Database(database).Collection(collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	if oid, ok := result["_id"].(primitive.ObjectID); ok {
		result["_id"] = oid.Hex()
	}
	return result, nil
}

//func (m *MongoDBClient) FindMany(ctx context.Context, database, collection string, filter interface{}) ([]bson.M, error) {
//	cursor, err := m.Client.Database(database).Collection(collection).Find(ctx, filter)
//	if err != nil {
//		return nil, err
//	}
//	defer func(cursor *mongo.Cursor, ctx context2.Context) {
//		err := cursor.Close(ctx)
//		if err != nil {
//
//		}
//	}(cursor, ctx)
//
//	var results []bson.M
//	for cursor.Next(ctx) {
//		var result bson.M
//		if err := cursor.Decode(&result); err != nil {
//			return nil, err
//		}
//		if oid, ok := result["_id"].(primitive.ObjectID); ok {
//			result["_id"] = oid.Hex()
//		}
//		results = append(results, result)
//	}
//	if err := cursor.Err(); err != nil {
//		return nil, err
//	}
//	return results, nil
//}

func (m *MongoDBClient) FindMany(ctx context.Context, database, collection string, filter interface{}) []bson.M {
	results := make([]bson.M, 0)
	cursor, err := m.Client.Database(database).Collection(collection).Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return results
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error decoding document: %v", err)
			return results
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return results
	}
	return results
}
func (m *MongoDBClient) UpdateOne(ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	result, err := m.Client.Database(database).Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoDBClient) DeleteOne(ctx context.Context, database, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := m.Client.Database(database).Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
