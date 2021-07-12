// Code generated by mockery v2.8.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	logger "github.com/smartcontractkit/chainlink/core/logger"

	mock "github.com/stretchr/testify/mock"

	pipeline "github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// Runner is an autogenerated mock type for the Runner type
type Runner struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Runner) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExecuteAndInsertFinishedRun provides a mock function with given fields: ctx, spec, vars, l, saveSuccessfulTaskRuns
func (_m *Runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (int64, pipeline.FinalResult, error) {
	ret := _m.Called(ctx, spec, vars, l, saveSuccessfulTaskRuns)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) int64); ok {
		r0 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 pipeline.FinalResult
	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) pipeline.FinalResult); ok {
		r1 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r1 = ret.Get(1).(pipeline.FinalResult)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger, bool) error); ok {
		r2 = rf(ctx, spec, vars, l, saveSuccessfulTaskRuns)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ExecuteRun provides a mock function with given fields: ctx, spec, vars, l
func (_m *Runner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (pipeline.Run, pipeline.TaskRunResults, error) {
	ret := _m.Called(ctx, spec, vars, l)

	var r0 pipeline.Run
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) pipeline.Run); ok {
		r0 = rf(ctx, spec, vars, l)
	} else {
		r0 = ret.Get(0).(pipeline.Run)
	}

	var r1 pipeline.TaskRunResults
	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) pipeline.TaskRunResults); ok {
		r1 = rf(ctx, spec, vars, l)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(pipeline.TaskRunResults)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, pipeline.Spec, pipeline.Vars, logger.Logger) error); ok {
		r2 = rf(ctx, spec, vars, l)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Healthy provides a mock function with given fields:
func (_m *Runner) Healthy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertFinishedRun provides a mock function with given fields: db, run, trrs, saveSuccessfulTaskRuns
func (_m *Runner) InsertFinishedRun(db *gorm.DB, run pipeline.Run, trrs pipeline.TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error) {
	ret := _m.Called(db, run, trrs, saveSuccessfulTaskRuns)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*gorm.DB, pipeline.Run, pipeline.TaskRunResults, bool) int64); ok {
		r0 = rf(db, run, trrs, saveSuccessfulTaskRuns)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, pipeline.Run, pipeline.TaskRunResults, bool) error); ok {
		r1 = rf(db, run, trrs, saveSuccessfulTaskRuns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ready provides a mock function with given fields:
func (_m *Runner) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: ctx, run, l, saveSuccessfulTaskRuns
func (_m *Runner) Run(ctx context.Context, run *pipeline.Run, l logger.Logger, saveSuccessfulTaskRuns bool) (bool, error) {
	ret := _m.Called(ctx, run, l, saveSuccessfulTaskRuns)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *pipeline.Run, logger.Logger, bool) bool); ok {
		r0 = rf(ctx, run, l, saveSuccessfulTaskRuns)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pipeline.Run, logger.Logger, bool) error); ok {
		r1 = rf(ctx, run, l, saveSuccessfulTaskRuns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields:
func (_m *Runner) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TestInsertFinishedRun provides a mock function with given fields: db, jobID, jobName, jobType, specID
func (_m *Runner) TestInsertFinishedRun(db *gorm.DB, jobID int32, jobName string, jobType string, specID int32) (int64, error) {
	ret := _m.Called(db, jobID, jobName, jobType, specID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*gorm.DB, int32, string, string, int32) int64); ok {
		r0 = rf(db, jobID, jobName, jobType, specID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, int32, string, string, int32) error); ok {
		r1 = rf(db, jobID, jobName, jobType, specID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
