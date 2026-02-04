package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/language-platform/user-service/internal/model"
	"github.com/language-platform/user-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*model.UserResponse, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.UserProfile, error)
}

type userService struct {
	repo  repository.UserRepository
	redis *redis.Client
}

func NewUserService(repo repository.UserRepository, redis *redis.Client) UserService {
	return &userService{
		repo:  repo,
		redis: redis,
	}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%s", id.String())
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user model.UserResponse
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			logrus.Debugf("User %s loaded from cache", id)
			return &user, nil
		}
	}

	// Load from database
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Profile:   user.Profile,
	}

	// Cache the result
	data, _ := json.Marshal(response)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return response, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	responses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = &model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Profile:   user.Profile,
		}
	}

	return responses, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", id.String())
	s.redis.Del(ctx, cacheKey)

	return s.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", id.String())
	s.redis.Del(ctx, cacheKey)

	return nil
}

func (s *userService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error) {
	profile, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return profile, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.UserProfile, error) {
	profile, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		// Create new profile if not exists
		profile = &model.UserProfile{
			UserID: userID,
		}
	}

	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		profile.PhoneNumber = *req.PhoneNumber
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Country != nil {
		profile.Country = *req.Country
	}
	if req.City != nil {
		profile.City = *req.City
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}
	if req.Language != nil {
		profile.Language = *req.Language
	}
	if req.Timezone != nil {
		profile.Timezone = *req.Timezone
	}

	var saveErr error
	if profile.ID == uuid.Nil {
		saveErr = s.repo.CreateProfile(ctx, profile)
	} else {
		saveErr = s.repo.UpdateProfile(ctx, profile)
	}

	if saveErr != nil {
		return nil, fmt.Errorf("failed to save user profile: %w", saveErr)
	}

	// Invalidate user cache
	cacheKey := fmt.Sprintf("user:%s", userID.String())
	s.redis.Del(ctx, cacheKey)

	return profile, nil
}
