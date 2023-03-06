package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/KarolosLykos/sportsnews/domain"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

const (
	sportsNewsDB       = "sportsnews"
	articlesCollection = "articles"
)

var (
	ErrGetByID = errors.New("repository: getByID")
	ErrList    = errors.New("repository: list")
	ErrUpsert  = errors.New("repository: upsert")
)

type mongoRepository struct {
	logger logger.Logger
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client, logger logger.Logger) *mongoRepository {
	return &mongoRepository{
		client: client,
		logger: logger,
	}
}

func (m *mongoRepository) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	article := &domain.Article{}

	opts := options.FindOne()
	if err := m.articlesCollection().FindOne(ctx, bson.M{"_id": id}, opts).Decode(article); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGetByID, err)
	}

	return article, nil
}

func (m *mongoRepository) List(ctx context.Context) (*domain.Articles, error) {
	count, err := m.articlesCollection().CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	if count == 0 {
		return &domain.Articles{Articles: make([]*domain.Article, 0), Total: 0}, nil
	}
	cursor, err := m.articlesCollection().Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	defer cursor.Close(ctx)

	articles := make([]*domain.Article, 0)
	for cursor.Next(ctx) {
		a := &domain.Article{}
		if err = cursor.Decode(a); err != nil {
			return nil, fmt.Errorf("%w:%v", ErrList, err)
		}
		articles = append(articles, a)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	return &domain.Articles{Total: count, Articles: articles}, nil
}

func (m *mongoRepository) Upsert(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	filter := bson.D{{Key: "articleID", Value: article.ArticleID}}
	update := bson.D{{Key: "$set", Value: article}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	if err := m.articlesCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(article); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrUpsert, err)
	}

	return article, nil
}

func (m *mongoRepository) articlesCollection() *mongo.Collection {
	return m.client.Database(sportsNewsDB).Collection(articlesCollection)
}
