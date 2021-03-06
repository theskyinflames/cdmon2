// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package repository

import (
	"github.com/theskyinflames/cdmon2/app/config"
	"sync"
)

var (
	lockStoreMockClose   sync.RWMutex
	lockStoreMockConnect sync.RWMutex
	lockStoreMockGet     sync.RWMutex
	lockStoreMockGetAll  sync.RWMutex
	lockStoreMockRemove  sync.RWMutex
	lockStoreMockSet     sync.RWMutex
)

// Ensure, that StoreMock does implement Store.
// If this is not the case, regenerate this file with moq.
var _ Store = &StoreMock{}

// StoreMock is a mock implementation of Store.
//
//     func TestSomethingThatUsesStore(t *testing.T) {
//
//         // make and configure a mocked Store
//         mockedStore := &StoreMock{
//             CloseFunc: func() error {
// 	               panic("mock out the Close method")
//             },
//             ConnectFunc: func() error {
// 	               panic("mock out the Connect method")
//             },
//             GetFunc: func(key string, item interface{}) (interface{}, error) {
// 	               panic("mock out the Get method")
//             },
//             GetAllFunc: func(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error) {
// 	               panic("mock out the GetAll method")
//             },
//             RemoveFunc: func(key string) error {
// 	               panic("mock out the Remove method")
//             },
//             SetFunc: func(key string, item interface{}) error {
// 	               panic("mock out the Set method")
//             },
//         }
//
//         // use mockedStore in code that requires Store
//         // and then make assertions.
//
//     }
type StoreMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// ConnectFunc mocks the Connect method.
	ConnectFunc func() error

	// GetFunc mocks the Get method.
	GetFunc func(key string, item interface{}) (interface{}, error)

	// GetAllFunc mocks the GetAll method.
	GetAllFunc func(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error)

	// RemoveFunc mocks the Remove method.
	RemoveFunc func(key string) error

	// SetFunc mocks the Set method.
	SetFunc func(key string, item interface{}) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Connect holds details about calls to the Connect method.
		Connect []struct {
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Key is the key argument value.
			Key string
			// Item is the item argument value.
			Item interface{}
		}
		// GetAll holds details about calls to the GetAll method.
		GetAll []struct {
			// Pattern is the pattern argument value.
			Pattern string
			// EmptyRecordFunc is the emptyRecordFunc argument value.
			EmptyRecordFunc config.EmptyRecordFunc
		}
		// Remove holds details about calls to the Remove method.
		Remove []struct {
			// Key is the key argument value.
			Key string
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// Key is the key argument value.
			Key string
			// Item is the item argument value.
			Item interface{}
		}
	}
}

// Close calls CloseFunc.
func (mock *StoreMock) Close() error {
	if mock.CloseFunc == nil {
		panic("StoreMock.CloseFunc: method is nil but Store.Close was just called")
	}
	callInfo := struct {
	}{}
	lockStoreMockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	lockStoreMockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedStore.CloseCalls())
func (mock *StoreMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	lockStoreMockClose.RLock()
	calls = mock.calls.Close
	lockStoreMockClose.RUnlock()
	return calls
}

// Connect calls ConnectFunc.
func (mock *StoreMock) Connect() error {
	if mock.ConnectFunc == nil {
		panic("StoreMock.ConnectFunc: method is nil but Store.Connect was just called")
	}
	callInfo := struct {
	}{}
	lockStoreMockConnect.Lock()
	mock.calls.Connect = append(mock.calls.Connect, callInfo)
	lockStoreMockConnect.Unlock()
	return mock.ConnectFunc()
}

// ConnectCalls gets all the calls that were made to Connect.
// Check the length with:
//     len(mockedStore.ConnectCalls())
func (mock *StoreMock) ConnectCalls() []struct {
} {
	var calls []struct {
	}
	lockStoreMockConnect.RLock()
	calls = mock.calls.Connect
	lockStoreMockConnect.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *StoreMock) Get(key string, item interface{}) (interface{}, error) {
	if mock.GetFunc == nil {
		panic("StoreMock.GetFunc: method is nil but Store.Get was just called")
	}
	callInfo := struct {
		Key  string
		Item interface{}
	}{
		Key:  key,
		Item: item,
	}
	lockStoreMockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	lockStoreMockGet.Unlock()
	return mock.GetFunc(key, item)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedStore.GetCalls())
func (mock *StoreMock) GetCalls() []struct {
	Key  string
	Item interface{}
} {
	var calls []struct {
		Key  string
		Item interface{}
	}
	lockStoreMockGet.RLock()
	calls = mock.calls.Get
	lockStoreMockGet.RUnlock()
	return calls
}

// GetAll calls GetAllFunc.
func (mock *StoreMock) GetAll(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error) {
	if mock.GetAllFunc == nil {
		panic("StoreMock.GetAllFunc: method is nil but Store.GetAll was just called")
	}
	callInfo := struct {
		Pattern         string
		EmptyRecordFunc config.EmptyRecordFunc
	}{
		Pattern:         pattern,
		EmptyRecordFunc: emptyRecordFunc,
	}
	lockStoreMockGetAll.Lock()
	mock.calls.GetAll = append(mock.calls.GetAll, callInfo)
	lockStoreMockGetAll.Unlock()
	return mock.GetAllFunc(pattern, emptyRecordFunc)
}

// GetAllCalls gets all the calls that were made to GetAll.
// Check the length with:
//     len(mockedStore.GetAllCalls())
func (mock *StoreMock) GetAllCalls() []struct {
	Pattern         string
	EmptyRecordFunc config.EmptyRecordFunc
} {
	var calls []struct {
		Pattern         string
		EmptyRecordFunc config.EmptyRecordFunc
	}
	lockStoreMockGetAll.RLock()
	calls = mock.calls.GetAll
	lockStoreMockGetAll.RUnlock()
	return calls
}

// Remove calls RemoveFunc.
func (mock *StoreMock) Remove(key string) error {
	if mock.RemoveFunc == nil {
		panic("StoreMock.RemoveFunc: method is nil but Store.Remove was just called")
	}
	callInfo := struct {
		Key string
	}{
		Key: key,
	}
	lockStoreMockRemove.Lock()
	mock.calls.Remove = append(mock.calls.Remove, callInfo)
	lockStoreMockRemove.Unlock()
	return mock.RemoveFunc(key)
}

// RemoveCalls gets all the calls that were made to Remove.
// Check the length with:
//     len(mockedStore.RemoveCalls())
func (mock *StoreMock) RemoveCalls() []struct {
	Key string
} {
	var calls []struct {
		Key string
	}
	lockStoreMockRemove.RLock()
	calls = mock.calls.Remove
	lockStoreMockRemove.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *StoreMock) Set(key string, item interface{}) error {
	if mock.SetFunc == nil {
		panic("StoreMock.SetFunc: method is nil but Store.Set was just called")
	}
	callInfo := struct {
		Key  string
		Item interface{}
	}{
		Key:  key,
		Item: item,
	}
	lockStoreMockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	lockStoreMockSet.Unlock()
	return mock.SetFunc(key, item)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//     len(mockedStore.SetCalls())
func (mock *StoreMock) SetCalls() []struct {
	Key  string
	Item interface{}
} {
	var calls []struct {
		Key  string
		Item interface{}
	}
	lockStoreMockSet.RLock()
	calls = mock.calls.Set
	lockStoreMockSet.RUnlock()
	return calls
}
