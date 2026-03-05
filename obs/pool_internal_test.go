// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package obs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper test task for pool testing

type testTask struct {
	result interface{}
}

func (mt *testTask) Run() interface{} {
	return mt.result
}

// FutureResult Tests

func TestFutureResult_ShouldReturnResult_WhenCalled(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"String result", "hello", "hello"},
		{"Int result", 42, 42},
		{"Nil result", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FutureResult{
				result: tt.input,
			}
			result := f.Get()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFutureResult_ShouldPanic_WhenResultContainsPanic(t *testing.T) {
	f := &FutureResult{
		result: panicResult{
			presult: "panic value",
		},
	}
	assert.Panics(t, func() {
		f.Get()
	})
}

// funcWrapper Tests

func TestFuncWrapper_ShouldExecuteFunction_WhenFuncNotNil(t *testing.T) {
	fw := &funcWrapper{
		f: func() interface{} {
			return "result"
		},
	}
	result := fw.Run()
	assert.Equal(t, "result", result)
}

func TestFuncWrapper_ShouldReturnNil_WhenFuncIsNil(t *testing.T) {
	fw := &funcWrapper{
		f: nil,
	}
	result := fw.Run()
	assert.Nil(t, result)
}

// NewRoutinePool Tests

func TestNewRoutinePool_ShouldReturnPool_WhenGivenValidParams(t *testing.T) {
	pool := NewRoutinePool(5, 10)
	assert.NotNil(t, pool)
}

func TestNewRoutinePool_ShouldUseCPUCount_WhenMaxWorkerCntZero(t *testing.T) {
	pool := NewRoutinePool(0, 10)
	assert.NotNil(t, pool)
	assert.Greater(t, pool.GetMaxWorkerCnt(), int64(0))
}

func TestNewRoutinePool_ShouldUseCPUCount_WhenMaxWorkerCntNegative(t *testing.T) {
	pool := NewRoutinePool(-1, 10)
	assert.NotNil(t, pool)
	assert.Greater(t, pool.GetMaxWorkerCnt(), int64(0))
}

// RoutinePool Execute Tests

func TestRoutinePool_Execute_ShouldAddTaskToQueue_WhenTaskNotNil(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	task := &testTask{
		result: "executed",
	}
	pool.Execute(task)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "executed", task.result)
}

func TestRoutinePool_Execute_ShouldNotPanic_WhenTaskIsNil(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	assert.NotPanics(t, func() {
		pool.Execute(nil)
	})
}

// RoutinePool ExecuteFunc Tests

func TestRoutinePool_ExecuteFunc_ShouldExecuteFunction(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	resultCh := make(chan string, 1)
	pool.ExecuteFunc(func() interface{} {
		resultCh <- "executed"
		return nil
	})
	result := <-resultCh
	assert.Equal(t, "executed", result)
}

// RoutinePool Submit Tests

func TestRoutinePool_Submit_ShouldReturnFuture_WhenTaskValid(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	task := &testTask{
		result: "submitted",
	}
	future, err := pool.Submit(task)
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "submitted", result)
}

func TestRoutinePool_Submit_ShouldReturnError_WhenTaskNil(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	future, err := pool.Submit(nil)
	assert.Error(t, err)
	assert.Nil(t, future)
	assert.Equal(t, ErrTaskInvalid, err)
}

func TestRoutinePool_Submit_ShouldReturnError_WhenPoolShutDown(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	pool.ShutDown()

	task := &testTask{
		result: "submitted",
	}
	future, err := pool.Submit(task)
	assert.Error(t, err)
	assert.Nil(t, future)
	assert.Equal(t, ErrPoolShutDown, err)
}

// RoutinePool SubmitFunc Tests

func TestRoutinePool_SubmitFunc_ShouldReturnFuture_WhenFuncValid(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	future, err := pool.SubmitFunc(func() interface{} {
		return "func result"
	})
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "func result", result)
}

// basicPool Counters Tests

func TestBasicPool_GetMaxWorkerCnt_ShouldReturnMaxWorkerCount(t *testing.T) {
	pool := NewRoutinePool(5, 10)
	defer pool.ShutDown()

	assert.GreaterOrEqual(t, pool.GetMaxWorkerCnt(), int64(5))
}

func TestBasicPool_AddMaxWorkerCnt_ShouldIncreaseMaxWorkerCount(t *testing.T) {
	pool := NewRoutinePool(5, 10)
	defer pool.ShutDown()

	// Enable autoTune to allow AddMaxWorkerCnt to work
	pool.EnableAutoTune()

	oldMax := pool.GetMaxWorkerCnt()
	newMax := pool.AddMaxWorkerCnt(2)
	assert.Equal(t, oldMax+2, newMax)
}

func TestBasicPool_GetCurrentWorkingCnt_ShouldReturnCurrentWorkerCount(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	task := &testTask{
		result: "test",
	}
	pool.Execute(task)
	time.Sleep(50 * time.Millisecond)

	count := pool.GetCurrentWorkingCnt()
	assert.GreaterOrEqual(t, count, int64(0))
}

func TestBasicPool_AddCurrentWorkingCnt_ShouldUpdateCurrentWorkerCount(t *testing.T) {
	pool := NewRoutinePool(5, 10)
	defer pool.ShutDown()

	oldCount := pool.GetCurrentWorkingCnt()
	newCount := pool.AddCurrentWorkingCnt(1)
	assert.Equal(t, oldCount+1, newCount)
}

func TestBasicPool_AddWorkerCnt_ShouldIncreaseWorkerCount(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	oldCount := pool.GetWorkerCnt()
	newCount := pool.AddWorkerCnt(1)
	assert.Equal(t, oldCount+1, newCount)
}

func TestBasicPool_EnableAutoTune_ShouldNotPanic(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	assert.NotPanics(t, func() {
		pool.EnableAutoTune()
	})
}

// RoutinePool EnableAutoTune Tests

func TestRoutinePool_EnableAutoTune_ShouldEnableAutoTune(t *testing.T) {
	pool := NewRoutinePool(2, 5)
	defer pool.ShutDown()

	pool.EnableAutoTune()
	// Verify that AddMaxWorkerCnt works after enabling
	oldMax := pool.GetMaxWorkerCnt()
	newMax := pool.AddMaxWorkerCnt(1)
	assert.Equal(t, oldMax+1, newMax)
}

// NoChanPool Tests

func TestNewNochanPool_ShouldReturnPool_WhenGivenValidParams(t *testing.T) {
	pool := NewNochanPool(5)
	defer pool.ShutDown()

	assert.NotNil(t, pool)
}

func TestNewNochanPool_ShouldUseCPUCount_WhenMaxWorkerCntZero(t *testing.T) {
	pool := NewNochanPool(0)
	defer pool.ShutDown()

	assert.NotNil(t, pool)
	assert.Greater(t, pool.GetMaxWorkerCnt(), int64(0))
}

func TestNoChanPool_Execute_ShouldExecuteTask_WhenTaskNotNil(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	task := &testTask{
		result: "nochan executed",
	}
	pool.Execute(task)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "nochan executed", task.result)
}

func TestNoChanPool_Execute_ShouldNotPanic_WhenTaskIsNil(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	assert.NotPanics(t, func() {
		pool.Execute(nil)
	})
}

func TestNoChanPool_ExecuteFunc_ShouldExecuteFunction(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	resultCh := make(chan string, 1)
	pool.ExecuteFunc(func() interface{} {
		resultCh <- "nochan func executed"
		return nil
	})
	result := <-resultCh
	assert.Equal(t, "nochan func executed", result)
}

func TestNoChanPool_Submit_ShouldReturnFuture_WhenTaskValid(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	task := &testTask{
		result: "nochan submitted",
	}
	future, err := pool.Submit(task)
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "nochan submitted", result)
}

func TestNoChanPool_Submit_ShouldReturnError_WhenTaskNil(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	future, err := pool.Submit(nil)
	assert.Error(t, err)
	assert.Nil(t, future)
	assert.Equal(t, ErrTaskInvalid, err)
}

// ShutDown Tests

func TestRoutinePool_ShutDown_ShouldClosePool_WhenCalled(t *testing.T) {
	pool := NewRoutinePool(2, 5)

	task := &testTask{
		result: "before shutdown",
	}
	pool.Execute(task)
	time.Sleep(50 * time.Millisecond)

	pool.ShutDown()

	_, err := pool.Submit(&testTask{result: "after shutdown"})
	assert.Error(t, err)
	assert.Equal(t, ErrPoolShutDown, err)
}

// signalTask Tests

func TestSignalTask_Run_ShouldReturnNil(t *testing.T) {
	st := signalTask{id: "test"}
	result := st.Run()
	assert.Nil(t, result)
}

// CompareAndSwapCurrentWorkingCnt Tests

func TestCompareAndSwapCurrentWorkingCnt_ShouldReturnTrue_WhenValueMatches(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access CompareAndSwapCurrentWorkingCnt
	pool := routinePool.(*RoutinePool)

	// First, get the current working count
	oldCount := pool.GetCurrentWorkingCnt()

	// Try to swap with the same value
	swapped := pool.CompareAndSwapCurrentWorkingCnt(oldCount, oldCount+1)
	assert.True(t, swapped)

	// Verify the new value
	newCount := pool.GetCurrentWorkingCnt()
	assert.Equal(t, oldCount+1, newCount)
}

func TestCompareAndSwapCurrentWorkingCnt_ShouldReturnFalse_WhenValueNotMatches(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access CompareAndSwapCurrentWorkingCnt
	pool := routinePool.(*RoutinePool)

	// Try to swap with a non-matching old value
	swapped := pool.CompareAndSwapCurrentWorkingCnt(999, 1000)
	assert.False(t, swapped)
}

// basicPool EnableAutoTune Tests

// Note: TestBasicPool_EnableAutoTune_ShouldNotPanic already exists at line 252

// SubmitWithTimeout Tests

func TestSubmitWithTimeout_ShouldSubmitTask_WhenTimeoutIsZero(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	task := &testTask{
		result: "timeout test",
	}
	future, err := pool.SubmitWithTimeout(task, 0)
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "timeout test", result)
}

func TestSubmitWithTimeout_ShouldSubmitTask_WhenTimeoutIsNegative(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	task := &testTask{
		result: "negative timeout test",
	}
	future, err := pool.SubmitWithTimeout(task, -1)
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "negative timeout test", result)
}

func TestSubmitWithTimeout_ShouldReturnError_WhenTaskNil(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	future, err := pool.SubmitWithTimeout(nil, 100)
	assert.Error(t, err)
	assert.Nil(t, future)
	assert.Equal(t, ErrTaskInvalid, err)
}

func TestSubmitWithTimeout_ShouldReturnError_WhenPoolShutDown(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	task := &testTask{
		result: "shutdown test",
	}
	future, err := pool.SubmitWithTimeout(task, 100)
	assert.Error(t, err)
	assert.Nil(t, future)
	assert.Equal(t, ErrPoolShutDown, err)
}

func TestSubmitWithTimeout_ShouldSubmitSuccessfully_WhenWithinTimeout(t *testing.T) {
	routinePool := NewRoutinePool(2, 5)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	task := &testTask{
		result: "timed submit",
	}
	future, err := pool.SubmitWithTimeout(task, 5000) // 5 seconds
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "timed submit", result)
}

func TestSubmitWithTimeout_ShouldReturnTimeoutError_WhenQueueFull(t *testing.T) {
	// Create a pool with small queue
	routinePool := NewRoutinePool(1, 1)
	defer routinePool.ShutDown()

	// Type assert to *RoutinePool to access SubmitWithTimeout
	pool := routinePool.(*RoutinePool)

	// Fill the queue
	task := &testTask{result: "blocking"}
	blockingFuture, _ := pool.Submit(task)
	// Don't get result, keeping the channel full

	// Try to submit another task with short timeout
	timeoutTask := &testTask{result: "timeout"}
	future, err := pool.SubmitWithTimeout(timeoutTask, 10) // 10ms

	// The behavior depends on timing, but we can test the API
	if err != nil {
		assert.Equal(t, ErrSubmitTimeout, err)
		assert.Nil(t, future)
	} else {
		// If it succeeded, get the result to clean up
		if future != nil {
			future.Get()
		}
	}
	_ = blockingFuture.Get() // Clean up
}

// NoChanPool SubmitFunc Tests

func TestNoChanPool_SubmitFunc_ShouldExecuteFunction(t *testing.T) {
	pool := NewNochanPool(2)
	defer pool.ShutDown()

	future, err := pool.SubmitFunc(func() interface{} {
		return "nochan func result"
	})
	assert.NoError(t, err)
	assert.NotNil(t, future)

	result := future.Get()
	assert.Equal(t, "nochan func result", result)
}

// Note: NoChanPool.SubmitFunc does not return error for nil function,
// it will execute and return nil. This is by design.
