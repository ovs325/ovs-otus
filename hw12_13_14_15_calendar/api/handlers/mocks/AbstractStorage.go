// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/common"

	mock "github.com/stretchr/testify/mock"

	time "time"

	types "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
)

// AbstractStorage is an autogenerated mock type for the AbstractStorage type
type AbstractStorage struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *AbstractStorage) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connect provides a mock function with given fields: ctx
func (_m *AbstractStorage) Connect(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Connect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateEvent provides a mock function with given fields: ctx, event
func (_m *AbstractStorage) CreateEvent(ctx context.Context, event *types.EventModel) (int64, error) {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for CreateEvent")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.EventModel) (int64, error)); ok {
		return rf(ctx, event)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.EventModel) int64); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.EventModel) error); ok {
		r1 = rf(ctx, event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DelEvent provides a mock function with given fields: ctx, id
func (_m *AbstractStorage) DelEvent(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DelEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEventsForTimeInterval provides a mock function with given fields: ctx, start, end, datePaginate
func (_m *AbstractStorage) GetEventsForTimeInterval(ctx context.Context, start time.Time, end time.Time, datePaginate common.Paginate) (types.QueryPage[types.EventModel], error) {
	ret := _m.Called(ctx, start, end, datePaginate)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsForTimeInterval")
	}

	var r0 types.QueryPage[types.EventModel]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time, common.Paginate) (types.QueryPage[types.EventModel], error)); ok {
		return rf(ctx, start, end, datePaginate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time, common.Paginate) types.QueryPage[types.EventModel]); ok {
		r0 = rf(ctx, start, end, datePaginate)
	} else {
		r0 = ret.Get(0).(types.QueryPage[types.EventModel])
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time, common.Paginate) error); ok {
		r1 = rf(ctx, start, end, datePaginate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEvent provides a mock function with given fields: ctx, event
func (_m *AbstractStorage) UpdateEvent(ctx context.Context, event *types.EventModel) error {
	ret := _m.Called(ctx, event)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.EventModel) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAbstractStorage creates a new instance of AbstractStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAbstractStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *AbstractStorage {
	mock := &AbstractStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
