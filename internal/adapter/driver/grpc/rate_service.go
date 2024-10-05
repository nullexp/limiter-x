package grpc

import (
	"context"

	ratev1 "github.com/nullexp/limiter-x/internal/adapter/driver/grpc/proto/rate/v1"
	driverService "github.com/nullexp/limiter-x/internal/port/driver/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiterService implements the gRPC RateLimiterServiceServer interface.
type RateLimiterService struct {
	ratev1.UnimplementedRateLimiterServiceServer
	service driverService.RateLimiter
}

// NewRateLimiterService creates a new instance of RateLimiterService.
func NewRateLimiterService(rateLimiterService driverService.RateLimiter) *RateLimiterService {
	return &RateLimiterService{service: rateLimiterService}
}

// CheckRateLimit implements the rate-limiting logic for the CheckRateLimit gRPC call.
func (rls *RateLimiterService) CheckRateLimit(ctx context.Context, request *ratev1.CheckRateLimitRequest) (*ratev1.CheckRateLimitResponse, error) {
    // Call the RateLimit method from the service
    allowed, err := rls.service.RateLimit(ctx, request.UserId, int(request.Limit))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to check rate limit: %v", err)
    }

    // Return the response
    return &ratev1.CheckRateLimitResponse{
        Allowed: allowed,
        Message: "Rate limit checked",
    }, nil
}

// GetUserRateLimit implements the GetUserRateLimit gRPC call.
func (rls *RateLimiterService) GetUserRateLimit(ctx context.Context, request *ratev1.GetUserRateLimitRequest) (*ratev1.GetUserRateLimitResponse, error) {
	// Call the GetUserRateLimit method from the service
	rateLimitModel, err := rls.service.GetUserRateLimit(ctx, request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user rate limit: %v", err)
	}

	// Return the response with the rate limit model
	return &ratev1.GetUserRateLimitResponse{
		UserId:    rateLimitModel.UserId,
		Limit:     int32(rateLimitModel.Limit),
		Remaining: int32(rateLimitModel.Remaining),
		Window:    rateLimitModel.Window,
	}, nil
}

// UpdateUserRateLimit implements the UpdateUserRateLimit gRPC call.
func (rls *RateLimiterService) UpdateUserRateLimit(ctx context.Context, request *ratev1.UpdateUserRateLimitRequest) (*ratev1.UpdateUserRateLimitResponse, error) {
	// Call the UpdateUserRateLimit method from the service
	err := rls.service.UpdateUserRateLimit(ctx, request.UserId, int(request.NewLimit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user rate limit: %v", err)
	}

	// Return the response confirming the update
	return &ratev1.UpdateUserRateLimitResponse{
		UserId:       request.UserId,
		UpdatedLimit: request.NewLimit,
		Message:      "User rate limit updated successfully",
	}, nil
}
