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

const usersCollection = "users"

type UserRepositoryImpl struct {
	db *mongodb.Client
}

func NewUserRepository(db *mongodb.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	collection := r.db.GetCollection(usersCollection)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := r.db.GetCollection(usersCollection)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	collection := r.db.GetCollection(usersCollection)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *UserRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	collection := r.db.GetCollection(usersCollection)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	collection := r.db.GetCollection(usersCollection)

	user.UpdatedAt = time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)

	return err
}

// UpdateVerificationStatus updates user verification status
func (r *UserRepositoryImpl) UpdateVerificationStatus(ctx context.Context, email string, isVerified bool) error {
	collection := r.db.GetCollection(usersCollection)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"is_verified": isVerified,
				"updated_at":  time.Now(),
			},
		},
	)

	return err
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	collection := r.db.GetCollection(usersCollection)

	now := time.Now()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"last_login_at": &now,
				"updated_at":    now,
			},
		},
	)

	return err
}

// CreateIndexes creates necessary indexes for the users collection
func (r *UserRepositoryImpl) CreateIndexes(ctx context.Context) error {
	collection := r.db.GetCollection(usersCollection)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
