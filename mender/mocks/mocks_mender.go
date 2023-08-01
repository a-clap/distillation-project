// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/a-clap/distillation-ota/pkg/mender (interfaces: Signer,Device,Downloader,Installer,Rebooter,LoadSaver,Callbacks)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	mender "github.com/a-clap/distillation-ota/pkg/mender"
	device "github.com/a-clap/distillation-ota/pkg/mender/device"
	gomock "github.com/golang/mock/gomock"
)

// MockSigner is a mock of Signer interface.
type MockSigner struct {
	ctrl     *gomock.Controller
	recorder *MockSignerMockRecorder
}

// MockSignerMockRecorder is the mock recorder for MockSigner.
type MockSignerMockRecorder struct {
	mock *MockSigner
}

// NewMockSigner creates a new mock instance.
func NewMockSigner(ctrl *gomock.Controller) *MockSigner {
	mock := &MockSigner{ctrl: ctrl}
	mock.recorder = &MockSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSigner) EXPECT() *MockSignerMockRecorder {
	return m.recorder
}

// PublicKeyPEM mocks base method.
func (m *MockSigner) PublicKeyPEM() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKeyPEM")
	ret0, _ := ret[0].(string)
	return ret0
}

// PublicKeyPEM indicates an expected call of PublicKeyPEM.
func (mr *MockSignerMockRecorder) PublicKeyPEM() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKeyPEM", reflect.TypeOf((*MockSigner)(nil).PublicKeyPEM))
}

// Sign mocks base method.
func (m *MockSigner) Sign(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockSignerMockRecorder) Sign(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockSigner)(nil).Sign), arg0)
}

// Verify mocks base method.
func (m *MockSigner) Verify(arg0, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockSignerMockRecorder) Verify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockSigner)(nil).Verify), arg0, arg1)
}

// MockDevice is a mock of Device interface.
type MockDevice struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceMockRecorder
}

// MockDeviceMockRecorder is the mock recorder for MockDevice.
type MockDeviceMockRecorder struct {
	mock *MockDevice
}

// NewMockDevice creates a new mock instance.
func NewMockDevice(ctrl *gomock.Controller) *MockDevice {
	mock := &MockDevice{ctrl: ctrl}
	mock.recorder = &MockDeviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDevice) EXPECT() *MockDeviceMockRecorder {
	return m.recorder
}

// Attributes mocks base method.
func (m *MockDevice) Attributes() ([]device.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes")
	ret0, _ := ret[0].([]device.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Attributes indicates an expected call of Attributes.
func (mr *MockDeviceMockRecorder) Attributes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockDevice)(nil).Attributes))
}

// ID mocks base method.
func (m *MockDevice) ID() ([]device.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].([]device.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ID indicates an expected call of ID.
func (mr *MockDeviceMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockDevice)(nil).ID))
}

// Info mocks base method.
func (m *MockDevice) Info() (device.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info")
	ret0, _ := ret[0].(device.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info.
func (mr *MockDeviceMockRecorder) Info() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockDevice)(nil).Info))
}

// MockDownloader is a mock of Downloader interface.
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockDownloaderMockRecorder
}

// MockDownloaderMockRecorder is the mock recorder for MockDownloader.
type MockDownloaderMockRecorder struct {
	mock *MockDownloader
}

// NewMockDownloader creates a new mock instance.
func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &MockDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDownloader) EXPECT() *MockDownloaderMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockDownloader) Download(arg0, arg1 string) (chan int, chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0, arg1)
	ret0, _ := ret[0].(chan int)
	ret1, _ := ret[1].(chan error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Download indicates an expected call of Download.
func (mr *MockDownloaderMockRecorder) Download(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDownloader)(nil).Download), arg0, arg1)
}

// MockInstaller is a mock of Installer interface.
type MockInstaller struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerMockRecorder
}

// MockInstallerMockRecorder is the mock recorder for MockInstaller.
type MockInstallerMockRecorder struct {
	mock *MockInstaller
}

// NewMockInstaller creates a new mock instance.
func NewMockInstaller(ctrl *gomock.Controller) *MockInstaller {
	mock := &MockInstaller{ctrl: ctrl}
	mock.recorder = &MockInstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstaller) EXPECT() *MockInstallerMockRecorder {
	return m.recorder
}

// Install mocks base method.
func (m *MockInstaller) Install(arg0 string) (chan int, chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install", arg0)
	ret0, _ := ret[0].(chan int)
	ret1, _ := ret[1].(chan error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Install indicates an expected call of Install.
func (mr *MockInstallerMockRecorder) Install(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockInstaller)(nil).Install), arg0)
}

// MockRebooter is a mock of Rebooter interface.
type MockRebooter struct {
	ctrl     *gomock.Controller
	recorder *MockRebooterMockRecorder
}

// MockRebooterMockRecorder is the mock recorder for MockRebooter.
type MockRebooterMockRecorder struct {
	mock *MockRebooter
}

// NewMockRebooter creates a new mock instance.
func NewMockRebooter(ctrl *gomock.Controller) *MockRebooter {
	mock := &MockRebooter{ctrl: ctrl}
	mock.recorder = &MockRebooterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRebooter) EXPECT() *MockRebooterMockRecorder {
	return m.recorder
}

// Reboot mocks base method.
func (m *MockRebooter) Reboot() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reboot indicates an expected call of Reboot.
func (mr *MockRebooterMockRecorder) Reboot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*MockRebooter)(nil).Reboot))
}

// MockLoadSaver is a mock of LoadSaver interface.
type MockLoadSaver struct {
	ctrl     *gomock.Controller
	recorder *MockLoadSaverMockRecorder
}

// MockLoadSaverMockRecorder is the mock recorder for MockLoadSaver.
type MockLoadSaverMockRecorder struct {
	mock *MockLoadSaver
}

// NewMockLoadSaver creates a new mock instance.
func NewMockLoadSaver(ctrl *gomock.Controller) *MockLoadSaver {
	mock := &MockLoadSaver{ctrl: ctrl}
	mock.recorder = &MockLoadSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoadSaver) EXPECT() *MockLoadSaverMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockLoadSaver) Load(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Load indicates an expected call of Load.
func (mr *MockLoadSaverMockRecorder) Load(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockLoadSaver)(nil).Load), arg0)
}

// Save mocks base method.
func (m *MockLoadSaver) Save(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockLoadSaverMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockLoadSaver)(nil).Save), arg0, arg1)
}

// MockCallbacks is a mock of Callbacks interface.
type MockCallbacks struct {
	ctrl     *gomock.Controller
	recorder *MockCallbacksMockRecorder
}

// MockCallbacksMockRecorder is the mock recorder for MockCallbacks.
type MockCallbacksMockRecorder struct {
	mock *MockCallbacks
}

// NewMockCallbacks creates a new mock instance.
func NewMockCallbacks(ctrl *gomock.Controller) *MockCallbacks {
	mock := &MockCallbacks{ctrl: ctrl}
	mock.recorder = &MockCallbacksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCallbacks) EXPECT() *MockCallbacksMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockCallbacks) Error(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0)
}

// Error indicates an expected call of Error.
func (mr *MockCallbacksMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockCallbacks)(nil).Error), arg0)
}

// NextState mocks base method.
func (m *MockCallbacks) NextState(arg0 mender.DeploymentStatus) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextState", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// NextState indicates an expected call of NextState.
func (mr *MockCallbacksMockRecorder) NextState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextState", reflect.TypeOf((*MockCallbacks)(nil).NextState), arg0)
}

// Update mocks base method.
func (m *MockCallbacks) Update(arg0 mender.DeploymentStatus, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0, arg1)
}

// Update indicates an expected call of Update.
func (mr *MockCallbacksMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCallbacks)(nil).Update), arg0, arg1)
}
