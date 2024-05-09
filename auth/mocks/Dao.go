// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	auth "github.com/auth"
	mock "github.com/stretchr/testify/mock"
)

// Dao is an autogenerated mock type for the Dao type
type Dao struct {
	mock.Mock
}

// GetUser provides a mock function with given fields: Username
func (_m *Dao) GetUser(Username string) (*auth.User, error) {
	ret := _m.Called(Username)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *auth.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*auth.User, error)); ok {
		return rf(Username)
	}
	if rf, ok := ret.Get(0).(func(string) *auth.User); ok {
		r0 = rf(Username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(Username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUser provides a mock function with given fields: _a0
func (_m *Dao) InsertUser(_a0 auth.User) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for InsertUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(auth.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDao creates a new instance of Dao. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDao(t interface {
	mock.TestingT
	Cleanup(func())
}) *Dao {
	mock := &Dao{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}