// Code generated by MockGen. DO NOT EDIT.
// Source: osservice (interfaces: Time,Store,Net,Wifi,Update,UpdateCallbacks)

// Package mocks is a generated GoMock package.
package mocks

import (
	mender "mender"
	osservice "osservice"
	wifi "osservice/pkg/wifi"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTime is a mock of Time interface.
type MockTime struct {
	ctrl     *gomock.Controller
	recorder *MockTimeMockRecorder
}

// MockTimeMockRecorder is the mock recorder for MockTime.
type MockTimeMockRecorder struct {
	mock *MockTime
}

// NewMockTime creates a new mock instance.
func NewMockTime(ctrl *gomock.Controller) *MockTime {
	mock := &MockTime{ctrl: ctrl}
	mock.recorder = &MockTimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTime) EXPECT() *MockTimeMockRecorder {
	return m.recorder
}

// NTP mocks base method.
func (m *MockTime) NTP() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NTP")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NTP indicates an expected call of NTP.
func (mr *MockTimeMockRecorder) NTP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NTP", reflect.TypeOf((*MockTime)(nil).NTP))
}

// Now mocks base method.
func (m *MockTime) Now() (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Now indicates an expected call of Now.
func (mr *MockTimeMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTime)(nil).Now))
}

// SetNTP mocks base method.
func (m *MockTime) SetNTP(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNTP", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNTP indicates an expected call of SetNTP.
func (mr *MockTimeMockRecorder) SetNTP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNTP", reflect.TypeOf((*MockTime)(nil).SetNTP), arg0)
}

// SetNow mocks base method.
func (m *MockTime) SetNow(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNow", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNow indicates an expected call of SetNow.
func (mr *MockTimeMockRecorder) SetNow(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNow", reflect.TypeOf((*MockTime)(nil).SetNow), arg0)
}

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockStore) Load(arg0 string) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockStoreMockRecorder) Load(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStore)(nil).Load), arg0)
}

// Save mocks base method.
func (m *MockStore) Save(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockStoreMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStore)(nil).Save), arg0, arg1)
}

// MockNet is a mock of Net interface.
type MockNet struct {
	ctrl     *gomock.Controller
	recorder *MockNetMockRecorder
}

// MockNetMockRecorder is the mock recorder for MockNet.
type MockNetMockRecorder struct {
	mock *MockNet
}

// NewMockNet creates a new mock instance.
func NewMockNet(ctrl *gomock.Controller) *MockNet {
	mock := &MockNet{ctrl: ctrl}
	mock.recorder = &MockNetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNet) EXPECT() *MockNetMockRecorder {
	return m.recorder
}

// ListInterfaces mocks base method.
func (m *MockNet) ListInterfaces() []osservice.NetInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInterfaces")
	ret0, _ := ret[0].([]osservice.NetInterface)
	return ret0
}

// ListInterfaces indicates an expected call of ListInterfaces.
func (mr *MockNetMockRecorder) ListInterfaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInterfaces", reflect.TypeOf((*MockNet)(nil).ListInterfaces))
}

// MockWifi is a mock of Wifi interface.
type MockWifi struct {
	ctrl     *gomock.Controller
	recorder *MockWifiMockRecorder
}

// MockWifiMockRecorder is the mock recorder for MockWifi.
type MockWifiMockRecorder struct {
	mock *MockWifi
}

// NewMockWifi creates a new mock instance.
func NewMockWifi(ctrl *gomock.Controller) *MockWifi {
	mock := &MockWifi{ctrl: ctrl}
	mock.recorder = &MockWifiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWifi) EXPECT() *MockWifiMockRecorder {
	return m.recorder
}

// APs mocks base method.
func (m *MockWifi) APs() ([]wifi.AP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APs")
	ret0, _ := ret[0].([]wifi.AP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APs indicates an expected call of APs.
func (mr *MockWifiMockRecorder) APs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APs", reflect.TypeOf((*MockWifi)(nil).APs))
}

// Connect mocks base method.
func (m *MockWifi) Connect(arg0 wifi.Network) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockWifiMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockWifi)(nil).Connect), arg0)
}

// Connected mocks base method.
func (m *MockWifi) Connected() (wifi.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connected")
	ret0, _ := ret[0].(wifi.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connected indicates an expected call of Connected.
func (mr *MockWifiMockRecorder) Connected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connected", reflect.TypeOf((*MockWifi)(nil).Connected))
}

// Disconnect mocks base method.
func (m *MockWifi) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockWifiMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockWifi)(nil).Disconnect))
}

// MockUpdate is a mock of Update interface.
type MockUpdate struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateMockRecorder
}

// MockUpdateMockRecorder is the mock recorder for MockUpdate.
type MockUpdateMockRecorder struct {
	mock *MockUpdate
}

// NewMockUpdate creates a new mock instance.
func NewMockUpdate(ctrl *gomock.Controller) *MockUpdate {
	mock := &MockUpdate{ctrl: ctrl}
	mock.recorder = &MockUpdateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdate) EXPECT() *MockUpdateMockRecorder {
	return m.recorder
}

// AvailableReleases mocks base method.
func (m *MockUpdate) AvailableReleases() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailableReleases")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailableReleases indicates an expected call of AvailableReleases.
func (mr *MockUpdateMockRecorder) AvailableReleases() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailableReleases", reflect.TypeOf((*MockUpdate)(nil).AvailableReleases))
}

// ContinueUpdate mocks base method.
func (m *MockUpdate) ContinueUpdate() (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContinueUpdate")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// ContinueUpdate indicates an expected call of ContinueUpdate.
func (mr *MockUpdateMockRecorder) ContinueUpdate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContinueUpdate", reflect.TypeOf((*MockUpdate)(nil).ContinueUpdate))
}

// PullReleases mocks base method.
func (m *MockUpdate) PullReleases() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullReleases")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PullReleases indicates an expected call of PullReleases.
func (mr *MockUpdateMockRecorder) PullReleases() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullReleases", reflect.TypeOf((*MockUpdate)(nil).PullReleases))
}

// StopUpdate mocks base method.
func (m *MockUpdate) StopUpdate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopUpdate")
	ret0, _ := ret[0].(error)
	return ret0
}

// StopUpdate indicates an expected call of StopUpdate.
func (mr *MockUpdateMockRecorder) StopUpdate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopUpdate", reflect.TypeOf((*MockUpdate)(nil).StopUpdate))
}

// Update mocks base method.
func (m *MockUpdate) Update(arg0 string, arg1 osservice.UpdateCallbacks) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUpdateMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUpdate)(nil).Update), arg0, arg1)
}

// MockUpdateCallbacks is a mock of UpdateCallbacks interface.
type MockUpdateCallbacks struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateCallbacksMockRecorder
}

// MockUpdateCallbacksMockRecorder is the mock recorder for MockUpdateCallbacks.
type MockUpdateCallbacksMockRecorder struct {
	mock *MockUpdateCallbacks
}

// NewMockUpdateCallbacks creates a new mock instance.
func NewMockUpdateCallbacks(ctrl *gomock.Controller) *MockUpdateCallbacks {
	mock := &MockUpdateCallbacks{ctrl: ctrl}
	mock.recorder = &MockUpdateCallbacksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdateCallbacks) EXPECT() *MockUpdateCallbacksMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockUpdateCallbacks) Error(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0)
}

// Error indicates an expected call of Error.
func (mr *MockUpdateCallbacksMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockUpdateCallbacks)(nil).Error), arg0)
}

// NextState mocks base method.
func (m *MockUpdateCallbacks) NextState(arg0 mender.DeploymentStatus) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextState", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// NextState indicates an expected call of NextState.
func (mr *MockUpdateCallbacksMockRecorder) NextState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextState", reflect.TypeOf((*MockUpdateCallbacks)(nil).NextState), arg0)
}

// Update mocks base method.
func (m *MockUpdateCallbacks) Update(arg0 mender.DeploymentStatus, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0, arg1)
}

// Update indicates an expected call of Update.
func (mr *MockUpdateCallbacksMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUpdateCallbacks)(nil).Update), arg0, arg1)
}