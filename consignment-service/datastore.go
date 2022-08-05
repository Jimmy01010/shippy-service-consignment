package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMongoClient 创建mongo客户端
func CreateMongoClient(ctx context.Context, uri string, retry int32) (*mongo.Client, error) {
	// 使用给定的连接字符串创建连接
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	// 检查它是否正确以及数据存储是否可用
	if err := conn.Ping(ctx, nil); err != nil {
		if retry >= 3 {
			return nil, err
		}
		retry = retry + 1
		time.Sleep(time.Second * 2)
		return CreateMongoClient(ctx, uri, retry)
	}

	return conn, err
}
