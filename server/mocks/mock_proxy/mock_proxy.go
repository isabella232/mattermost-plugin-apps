// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/proxy (interfaces: Service)

// Package mock_proxy is a generated GoMock package.
package mock_proxy

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	apps "github.com/mattermost/mattermost-plugin-apps/apps"
	appclient "github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	config "github.com/mattermost/mattermost-plugin-apps/server/config"
	incoming "github.com/mattermost/mattermost-plugin-apps/server/incoming"
	proxy "github.com/mattermost/mattermost-plugin-apps/server/proxy"
	upstream "github.com/mattermost/mattermost-plugin-apps/upstream"
	utils "github.com/mattermost/mattermost-plugin-apps/utils"
	model "github.com/mattermost/mattermost-server/v6/model"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddBuiltinUpstream mocks base method.
func (m *MockService) AddBuiltinUpstream(arg0 apps.AppID, arg1 upstream.Upstream) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddBuiltinUpstream", arg0, arg1)
}

// AddBuiltinUpstream indicates an expected call of AddBuiltinUpstream.
func (mr *MockServiceMockRecorder) AddBuiltinUpstream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBuiltinUpstream", reflect.TypeOf((*MockService)(nil).AddBuiltinUpstream), arg0, arg1)
}

// Call mocks base method.
func (m *MockService) Call(arg0 *incoming.Request, arg1 apps.CallRequest) proxy.CallResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", arg0, arg1)
	ret0, _ := ret[0].(proxy.CallResponse)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockServiceMockRecorder) Call(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockService)(nil).Call), arg0, arg1)
}

// CanDeploy mocks base method.
func (m *MockService) CanDeploy(arg0 apps.DeployType) (bool, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanDeploy", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CanDeploy indicates an expected call of CanDeploy.
func (mr *MockServiceMockRecorder) CanDeploy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanDeploy", reflect.TypeOf((*MockService)(nil).CanDeploy), arg0)
}

// CompleteRemoteOAuth2 mocks base method.
func (m *MockService) CompleteRemoteOAuth2(arg0 *incoming.Request, arg1 apps.AppID, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteRemoteOAuth2", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteRemoteOAuth2 indicates an expected call of CompleteRemoteOAuth2.
func (mr *MockServiceMockRecorder) CompleteRemoteOAuth2(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteRemoteOAuth2", reflect.TypeOf((*MockService)(nil).CompleteRemoteOAuth2), arg0, arg1, arg2)
}

// Configure mocks base method.
func (m *MockService) Configure(arg0 config.Config, arg1 utils.Logger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure.
func (mr *MockServiceMockRecorder) Configure(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockService)(nil).Configure), arg0, arg1)
}

// DisableApp mocks base method.
func (m *MockService) DisableApp(arg0 *incoming.Request, arg1 apps.Context, arg2 apps.AppID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableApp", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableApp indicates an expected call of DisableApp.
func (mr *MockServiceMockRecorder) DisableApp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableApp", reflect.TypeOf((*MockService)(nil).DisableApp), arg0, arg1, arg2)
}

// EnableApp mocks base method.
func (m *MockService) EnableApp(arg0 *incoming.Request, arg1 apps.Context, arg2 apps.AppID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableApp", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableApp indicates an expected call of EnableApp.
func (mr *MockServiceMockRecorder) EnableApp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableApp", reflect.TypeOf((*MockService)(nil).EnableApp), arg0, arg1, arg2)
}

// GetAppBindings mocks base method.
func (m *MockService) GetAppBindings(arg0 *incoming.Request, arg1 apps.Context, arg2 apps.App) ([]apps.Binding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppBindings", arg0, arg1, arg2)
	ret0, _ := ret[0].([]apps.Binding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAppBindings indicates an expected call of GetAppBindings.
func (mr *MockServiceMockRecorder) GetAppBindings(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppBindings", reflect.TypeOf((*MockService)(nil).GetAppBindings), arg0, arg1, arg2)
}

// GetBindings mocks base method.
func (m *MockService) GetBindings(arg0 *incoming.Request, arg1 apps.Context) ([]apps.Binding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBindings", arg0, arg1)
	ret0, _ := ret[0].([]apps.Binding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBindings indicates an expected call of GetBindings.
func (mr *MockServiceMockRecorder) GetBindings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBindings", reflect.TypeOf((*MockService)(nil).GetBindings), arg0, arg1)
}

// GetInstalledApp mocks base method.
func (m *MockService) GetInstalledApp(arg0 *incoming.Request, arg1 apps.AppID) (*apps.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstalledApp", arg0, arg1)
	ret0, _ := ret[0].(*apps.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstalledApp indicates an expected call of GetInstalledApp.
func (mr *MockServiceMockRecorder) GetInstalledApp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstalledApp", reflect.TypeOf((*MockService)(nil).GetInstalledApp), arg0, arg1)
}

// GetInstalledApps mocks base method.
func (m *MockService) GetInstalledApps(arg0 *incoming.Request) []apps.App {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstalledApps", arg0)
	ret0, _ := ret[0].([]apps.App)
	return ret0
}

// GetInstalledApps indicates an expected call of GetInstalledApps.
func (mr *MockServiceMockRecorder) GetInstalledApps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstalledApps", reflect.TypeOf((*MockService)(nil).GetInstalledApps), arg0)
}

// GetListedApps mocks base method.
func (m *MockService) GetListedApps(arg0 *incoming.Request, arg1 string, arg2 bool) []apps.ListedApp {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListedApps", arg0, arg1, arg2)
	ret0, _ := ret[0].([]apps.ListedApp)
	return ret0
}

// GetListedApps indicates an expected call of GetListedApps.
func (mr *MockServiceMockRecorder) GetListedApps(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListedApps", reflect.TypeOf((*MockService)(nil).GetListedApps), arg0, arg1, arg2)
}

// GetManifest mocks base method.
func (m *MockService) GetManifest(arg0 *incoming.Request, arg1 apps.AppID) (*apps.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManifest", arg0, arg1)
	ret0, _ := ret[0].(*apps.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManifest indicates an expected call of GetManifest.
func (mr *MockServiceMockRecorder) GetManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManifest", reflect.TypeOf((*MockService)(nil).GetManifest), arg0, arg1)
}

// GetRemoteOAuth2ConnectURL mocks base method.
func (m *MockService) GetRemoteOAuth2ConnectURL(arg0 *incoming.Request, arg1 apps.AppID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemoteOAuth2ConnectURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemoteOAuth2ConnectURL indicates an expected call of GetRemoteOAuth2ConnectURL.
func (mr *MockServiceMockRecorder) GetRemoteOAuth2ConnectURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemoteOAuth2ConnectURL", reflect.TypeOf((*MockService)(nil).GetRemoteOAuth2ConnectURL), arg0, arg1)
}

// GetStatic mocks base method.
func (m *MockService) GetStatic(arg0 *incoming.Request, arg1 apps.AppID, arg2 string) (io.ReadCloser, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatic", arg0, arg1, arg2)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetStatic indicates an expected call of GetStatic.
func (mr *MockServiceMockRecorder) GetStatic(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatic", reflect.TypeOf((*MockService)(nil).GetStatic), arg0, arg1, arg2)
}

// InstallApp mocks base method.
func (m *MockService) InstallApp(arg0 *incoming.Request, arg1 apps.Context, arg2 apps.AppID, arg3 apps.DeployType, arg4 bool, arg5 string) (*apps.App, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallApp", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*apps.App)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InstallApp indicates an expected call of InstallApp.
func (mr *MockServiceMockRecorder) InstallApp(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallApp", reflect.TypeOf((*MockService)(nil).InstallApp), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Notify mocks base method.
func (m *MockService) Notify(arg0 apps.Context, arg1 apps.Subject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Notify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Notify indicates an expected call of Notify.
func (mr *MockServiceMockRecorder) Notify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockService)(nil).Notify), arg0, arg1)
}

// NotifyMessageHasBeenPosted mocks base method.
func (m *MockService) NotifyMessageHasBeenPosted(arg0 *model.Post, arg1 apps.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyMessageHasBeenPosted", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyMessageHasBeenPosted indicates an expected call of NotifyMessageHasBeenPosted.
func (mr *MockServiceMockRecorder) NotifyMessageHasBeenPosted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyMessageHasBeenPosted", reflect.TypeOf((*MockService)(nil).NotifyMessageHasBeenPosted), arg0, arg1)
}

// NotifyRemoteWebhook mocks base method.
func (m *MockService) NotifyRemoteWebhook(arg0 *incoming.Request, arg1 apps.AppID, arg2 apps.HTTPCallRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyRemoteWebhook", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyRemoteWebhook indicates an expected call of NotifyRemoteWebhook.
func (mr *MockServiceMockRecorder) NotifyRemoteWebhook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyRemoteWebhook", reflect.TypeOf((*MockService)(nil).NotifyRemoteWebhook), arg0, arg1, arg2)
}

// NotifyUserHasJoinedChannel mocks base method.
func (m *MockService) NotifyUserHasJoinedChannel(arg0 apps.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyUserHasJoinedChannel", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyUserHasJoinedChannel indicates an expected call of NotifyUserHasJoinedChannel.
func (mr *MockServiceMockRecorder) NotifyUserHasJoinedChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyUserHasJoinedChannel", reflect.TypeOf((*MockService)(nil).NotifyUserHasJoinedChannel), arg0)
}

// NotifyUserHasJoinedTeam mocks base method.
func (m *MockService) NotifyUserHasJoinedTeam(arg0 apps.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyUserHasJoinedTeam", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyUserHasJoinedTeam indicates an expected call of NotifyUserHasJoinedTeam.
func (mr *MockServiceMockRecorder) NotifyUserHasJoinedTeam(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyUserHasJoinedTeam", reflect.TypeOf((*MockService)(nil).NotifyUserHasJoinedTeam), arg0)
}

// NotifyUserHasLeftChannel mocks base method.
func (m *MockService) NotifyUserHasLeftChannel(arg0 apps.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyUserHasLeftChannel", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyUserHasLeftChannel indicates an expected call of NotifyUserHasLeftChannel.
func (mr *MockServiceMockRecorder) NotifyUserHasLeftChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyUserHasLeftChannel", reflect.TypeOf((*MockService)(nil).NotifyUserHasLeftChannel), arg0)
}

// NotifyUserHasLeftTeam mocks base method.
func (m *MockService) NotifyUserHasLeftTeam(arg0 apps.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyUserHasLeftTeam", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyUserHasLeftTeam indicates an expected call of NotifyUserHasLeftTeam.
func (mr *MockServiceMockRecorder) NotifyUserHasLeftTeam(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyUserHasLeftTeam", reflect.TypeOf((*MockService)(nil).NotifyUserHasLeftTeam), arg0)
}

// SynchronizeInstalledApps mocks base method.
func (m *MockService) SynchronizeInstalledApps() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SynchronizeInstalledApps")
	ret0, _ := ret[0].(error)
	return ret0
}

// SynchronizeInstalledApps indicates an expected call of SynchronizeInstalledApps.
func (mr *MockServiceMockRecorder) SynchronizeInstalledApps() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SynchronizeInstalledApps", reflect.TypeOf((*MockService)(nil).SynchronizeInstalledApps))
}

// UninstallApp mocks base method.
func (m *MockService) UninstallApp(arg0 *incoming.Request, arg1 apps.Context, arg2 apps.AppID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UninstallApp", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UninstallApp indicates an expected call of UninstallApp.
func (mr *MockServiceMockRecorder) UninstallApp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UninstallApp", reflect.TypeOf((*MockService)(nil).UninstallApp), arg0, arg1, arg2)
}

// UpdateAppListing mocks base method.
func (m *MockService) UpdateAppListing(arg0 *incoming.Request, arg1 appclient.UpdateAppListingRequest) (*apps.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAppListing", arg0, arg1)
	ret0, _ := ret[0].(*apps.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAppListing indicates an expected call of UpdateAppListing.
func (mr *MockServiceMockRecorder) UpdateAppListing(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAppListing", reflect.TypeOf((*MockService)(nil).UpdateAppListing), arg0, arg1)
}
