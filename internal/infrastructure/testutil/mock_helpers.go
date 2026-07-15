package testutil

import (
	"github.com/stretchr/testify/mock"
)

// AnyContext returns a mock matcher for context.Context
func AnyContext() interface{} {
	return mock.AnythingOfType("*context.emptyCtx")
}

// AnyUUID returns a mock matcher for uuid.UUID
func AnyUUID() interface{} {
	return mock.AnythingOfType("uuid.UUID")
}

// AnyString returns a mock matcher for string
func AnyString() interface{} {
	return mock.AnythingOfType("string")
}

// AnyInt returns a mock matcher for int
func AnyInt() interface{} {
	return mock.AnythingOfType("int")
}

// AnyError returns a mock matcher for error
func AnyError() interface{} {
	return mock.AnythingOfType("*errors.errorString")
}

// MatchFunc creates a custom mock matcher function
func MatchFunc(fn func(interface{}) bool) interface{} {
	return mock.MatchedBy(fn)
}

// MockCallBuilder helps build mock expectations fluently
type MockCallBuilder struct {
	mockObj *mock.Mock
	call    *mock.Call
}

// NewMockCallBuilder creates a new mock call builder
func NewMockCallBuilder(mockObj *mock.Mock, method string, args ...interface{}) *MockCallBuilder {
	call := mockObj.On(method, args...)
	return &MockCallBuilder{
		mockObj: mockObj,
		call:    call,
	}
}

// Return sets the return values
func (b *MockCallBuilder) Return(returnArgs ...interface{}) *MockCallBuilder {
	b.call.Return(returnArgs...)
	return b
}

// Once sets the expectation to be called once
func (b *MockCallBuilder) Once() *MockCallBuilder {
	b.call.Once()
	return b
}

// Twice sets the expectation to be called twice
func (b *MockCallBuilder) Twice() *MockCallBuilder {
	b.call.Twice()
	return b
}

// Times sets the expectation to be called n times
func (b *MockCallBuilder) Times(n int) *MockCallBuilder {
	b.call.Times(n)
	return b
}

// Maybe marks the call as optional
func (b *MockCallBuilder) Maybe() *MockCallBuilder {
	b.call.Maybe()
	return b
}

// Run sets a function to run when the method is called
func (b *MockCallBuilder) Run(fn func(args mock.Arguments)) *MockCallBuilder {
	b.call.Run(fn)
	return b
}
