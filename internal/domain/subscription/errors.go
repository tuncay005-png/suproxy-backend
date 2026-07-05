package subscription

import "errors"

var (
	ErrSubscriptionNotFound         = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists    = errors.New("subscription already exists")
	ErrSubscriptionAlreadyActive    = errors.New("subscription already active")
	ErrSubscriptionAlreadySuspended = errors.New("subscription already suspended")
	ErrSubscriptionAlreadyCancelled = errors.New("subscription already cancelled")
	ErrSubscriptionAlreadyExpired   = errors.New("subscription already expired")
	ErrInvalidUserID                = errors.New("invalid user id")
	ErrInvalidPlan                  = errors.New("invalid plan")
	ErrInvalidPlanName              = errors.New("invalid plan name")
	ErrInvalidTrafficLimit          = errors.New("invalid traffic limit")
	ErrInvalidDeviceLimit           = errors.New("invalid device limit")
	ErrPlanNotFound                 = errors.New("plan not found")
	ErrPlanAlreadyExists            = errors.New("plan already exists")
)
