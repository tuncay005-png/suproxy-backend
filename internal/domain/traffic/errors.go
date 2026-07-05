package traffic

import "errors"

var (
	ErrUsageNotFound         = errors.New("usage not found")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidSubscriptionID = errors.New("invalid subscription id")
	ErrInvalidNodeID         = errors.New("invalid node id")
	ErrInvalidBytes          = errors.New("invalid bytes value")
)
