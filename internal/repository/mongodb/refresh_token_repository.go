package repository

import (
	"context"
	"errors"
	"time"

	domain "github.com/muzammil-cyber/gin-erp/internal/domain/auth"
	"github.com/muzammil-cyber/gin-erp/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const refreshTokensCollection = "refresh_tokens"

type RefreshTokenRepositoryImpl struct {
	db *mongodb.Client
}

func NewRefreshTokenRepository(db *mongodb.Client) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{
		db: db,
	}
}

// Create creates a new refresh token
func (r *RefreshTokenRepositoryImpl) Create(ctx context.Context, token *domain.RefreshToken) error {
	collection := r.db.GetCollection(refreshTokensCollection)

	token.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	token.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByToken finds a refresh token by token string
func (r *RefreshTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	collection := r.db.GetCollection(refreshTokensCollection)

	var refreshToken domain.RefreshToken
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&refreshToken)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	return &refreshToken, nil
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepositoryImpl) Revoke(ctx context.Context, token string) error {
	collection := r.db.GetCollection(refreshTokensCollection)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)

	return err
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepositoryImpl) RevokeAllForUser(ctx context.Context, userID primitive.ObjectID) error {
	collection := r.db.GetCollection(refreshTokensCollection)

	_, err := collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "is_revoked": false},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)

	return err
}

// CreateIndexes creates necessary indexes
func (r *RefreshTokenRepositoryImpl) CreateIndexes(ctx context.Context) error {
	collection := r.db.GetCollection(refreshTokensCollection)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
