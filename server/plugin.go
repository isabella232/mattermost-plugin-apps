// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"bytes"
	gohttp "net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-api/cluster"
	mmtelemetry "github.com/mattermost/mattermost-plugin-api/experimental/telemetry"
	"github.com/mattermost/mattermost-plugin-api/i18n"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/appservices"
	"github.com/mattermost/mattermost-plugin-apps/server/builtin"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/httpin"
	"github.com/mattermost/mattermost-plugin-apps/server/httpin/gateway"
	"github.com/mattermost/mattermost-plugin-apps/server/httpin/restapi"
	"github.com/mattermost/mattermost-plugin-apps/server/httpout"
	"github.com/mattermost/mattermost-plugin-apps/server/proxy"
	"github.com/mattermost/mattermost-plugin-apps/server/session"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/server/telemetry"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

type Plugin struct {
	plugin.MattermostPlugin
	manifest model.Manifest

	conf config.Service
	log  utils.Logger

	store          *store.Service
	appservices    appservices.Service
	proxy          proxy.Service
	sessionService session.Service

	httpIn  httpin.Service
	httpOut httpout.Service

	telemetryClient mmtelemetry.Client
	tracker         *telemetry.Telemetry
}

func NewPlugin(pluginManifest model.Manifest) *Plugin {
	return &Plugin{
		manifest: pluginManifest,
	}
}

var infoTemplate = template.Must(template.New("info").Parse(
	"Version: {{.Version}}, {{.URL}}, built {{.BuildDate}}, Cloud Mode: {{.CloudMode}}, Developer Mode: {{.DeveloperMode}}, Allow install over HTTP: {{.AllowHTTPApps}}"))

func (p *Plugin) OnActivate() (err error) {
	mm := pluginapi.NewClient(p.API, p.Driver)
	p.log = utils.NewPluginLogger(mm)

	botUserID, err := mm.Bot.EnsureBot(&model.Bot{
		Username:    config.BotUsername,
		DisplayName: config.BotDisplayName,
		Description: config.BotDescription,
	}, pluginapi.ProfileImagePath("assets/profile.png"))
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot account")
	}

	i18nBundle, err := i18n.InitBundle(p.API, filepath.Join("assets", "i18n"))
	if err != nil {
		return err
	}

	p.telemetryClient, err = mmtelemetry.NewRudderClient()
	if err != nil {
		p.API.LogWarn("telemetry client not started", "error", err.Error())
	}

	p.tracker = telemetry.NewTelemetry(nil)

	p.conf = config.NewService(mm, p.manifest, botUserID, p.tracker, i18nBundle)
	stored := config.StoredConfig{}
	_ = mm.Configuration.LoadPluginConfiguration(&stored)
	err = p.conf.Reconfigure(stored)
	if err != nil {
		return errors.Wrap(err, "failed to load initial configuration")
	}
	conf := p.conf.Get()
	infoBuf := &bytes.Buffer{}
	_ = infoTemplate.Execute(infoBuf, conf.InfoTemplateData())
	p.log.Debugf("OnActivate: %s", infoBuf.String())

	p.httpOut = httpout.NewService(p.conf)

	p.store, err = store.MakeService(p.log, p.conf, p.httpOut)
	if err != nil {
		return errors.Wrap(err, "failed to initialize persistent store")
	}
	p.store.App.InitBuiltin(builtin.App(conf))
	p.log.Debugf("Initialized persistent store")

	p.appservices = appservices.NewService(p.conf, p.store)

	p.sessionService = session.NewService(mm, p.store)

	mutex, err := cluster.NewMutex(p.API, config.KVClusterMutexKey)
	if err != nil {
		return errors.Wrapf(err, "failed creating cluster mutex")
	}

	p.proxy = proxy.NewService(p.conf, p.store, mutex, p.httpOut, p.sessionService, p.appservices)
	err = p.proxy.Configure(conf, p.log)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize app proxy service")
	}
	p.proxy.AddBuiltinUpstream(
		builtin.AppID,
		builtin.NewBuiltinApp(p.conf, p.proxy, p.appservices, p.httpOut, p.sessionService),
	)
	p.log.Debugf("Initialized the app proxy")

	p.httpIn = httpin.NewService(mm, mux.NewRouter(), p.conf, p.log, p.sessionService, p.proxy, p.appservices,
		restapi.Init,
		gateway.Init,
	)

	if conf.MattermostCloudMode {
		err = p.proxy.SynchronizeInstalledApps()
		if err != nil {
			p.log.WithError(err).Errorf("Failed to synchronize apps metadata")
		} else {
			p.log.Debugf("Synchronized the installed apps metadata")
		}
	}
	p.log.Infof("Plugin activated")

	p.conf.MattermostAPI().Frontend.PublishWebSocketEvent(config.WebSocketEventPluginEnabled, conf.GetPluginVersionInfo(), &model.WebsocketBroadcast{})

	return nil
}

func (p *Plugin) OnDeactivate() error {
	conf := p.conf.Get()
	p.conf.MattermostAPI().Frontend.PublishWebSocketEvent(config.WebSocketEventPluginDisabled, conf.GetPluginVersionInfo(), &model.WebsocketBroadcast{})

	if p.telemetryClient != nil {
		err := p.telemetryClient.Close()
		if err != nil {
			p.API.LogWarn("OnDeactivate: failed to close telemetryClient", "error", err.Error())
		}
	}

	return nil
}

func (p *Plugin) OnConfigurationChange() (err error) {
	defer func() {
		if err != nil {
			p.log.WithError(err).Errorf("Failed to reconfigure")
		}
	}()

	if p.conf == nil {
		// pre-activate, nothing to do.
		return nil
	}

	enableDiagnostics := false
	if config := p.API.GetConfig(); config != nil {
		if configValue := config.LogSettings.EnableDiagnostics; configValue != nil {
			enableDiagnostics = *configValue
		}
	}
	updatedTracker := mmtelemetry.NewTracker(p.telemetryClient, p.API.GetDiagnosticId(), p.API.GetServerVersion(), manifest.Id, manifest.Version, "appsFramework", enableDiagnostics)
	p.tracker.UpdateTracker(updatedTracker)

	mm := pluginapi.NewClient(p.API, p.Driver)
	stored := config.StoredConfig{}
	_ = mm.Configuration.LoadPluginConfiguration(&stored)

	return p.conf.Reconfigure(stored, p.store.App, p.store.Manifest, p.proxy)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w gohttp.ResponseWriter, req *gohttp.Request) {
	p.httpIn.ServeHTTP(c, w, req)
}

func (p *Plugin) UserHasBeenCreated(_ *plugin.Context, user *model.User) {
	err := p.proxy.Notify(
		apps.Context{
			UserID: user.Id,
			ExpandedContext: apps.ExpandedContext{
				User: user,
			},
		},
		apps.SubjectUserCreated)
	if err != nil {
		p.log.WithError(err).Debugf("Error handling UserHasBeenCreated")
	}
}

func (p *Plugin) UserHasJoinedChannel(_ *plugin.Context, cm *model.ChannelMember, actingUser *model.User) {
	err := p.proxy.NotifyUserHasJoinedChannel(p.newChannelMemberContext(cm))
	if err != nil {
		p.log.WithError(err).Debugf("Error handling UserHasJoinedChannel")
	}
}

func (p *Plugin) UserHasLeftChannel(_ *plugin.Context, cm *model.ChannelMember, actingUser *model.User) {
	err := p.proxy.NotifyUserHasLeftChannel(p.newChannelMemberContext(cm))
	if err != nil {
		p.log.WithError(err).Debugf("Error handling UserHasLeftChannel")
	}
}

func (p *Plugin) UserHasJoinedTeam(_ *plugin.Context, tm *model.TeamMember, actingUser *model.User) {
	err := p.proxy.NotifyUserHasJoinedTeam(p.newTeamMemberContext(tm))
	if err != nil {
		p.log.WithError(err).Debugf("Error handling UserHasJoinedTeam")
	}
}

func (p *Plugin) UserHasLeftTeam(_ *plugin.Context, tm *model.TeamMember, actingUser *model.User) {
	err := p.proxy.NotifyUserHasLeftTeam(p.newTeamMemberContext(tm))
	if err != nil {
		p.log.WithError(err).Debugf("Error handling UserHasLeftTeam")
	}
}

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	shouldProcessMessage, err := p.conf.MattermostAPI().Post.ShouldProcessMessage(post, pluginapi.BotID(p.conf.Get().BotUserID))
	if err != nil {
		p.log.WithError(err).Errorf("Error while checking if the message should be processed")
		return
	}

	if !shouldProcessMessage {
		return
	}

	err = p.proxy.NotifyMessageHasBeenPosted(post, apps.Context{
		UserAgentContext: apps.UserAgentContext{
			PostID:     post.Id,
			RootPostID: post.RootId,
			ChannelID:  post.ChannelId,
		},
		UserID: post.UserId,
		ExpandedContext: apps.ExpandedContext{
			Post: post,
		},
	})
	if err != nil {
		p.log.WithError(err).Debugf("Error handling MessageHasBeenPosted")
	}
}

func (p *Plugin) ChannelHasBeenCreated(_ *plugin.Context, ch *model.Channel) {
	err := p.proxy.Notify(
		apps.Context{
			UserAgentContext: apps.UserAgentContext{
				TeamID:    ch.TeamId,
				ChannelID: ch.Id,
			},
			UserID: ch.CreatorId,
			ExpandedContext: apps.ExpandedContext{
				Channel: ch,
			},
		},
		apps.SubjectChannelCreated)
	if err != nil {
		p.log.WithError(err).Debugf("Error handling ChannelHasBeenCreated")
	}
}

func (p *Plugin) newTeamMemberContext(tm *model.TeamMember) apps.Context {
	return apps.Context{
		UserAgentContext: apps.UserAgentContext{
			TeamID: tm.TeamId,
		},
		UserID: tm.UserId,
	}
}

func (p *Plugin) newChannelMemberContext(cm *model.ChannelMember) apps.Context {
	return apps.Context{
		UserAgentContext: apps.UserAgentContext{
			ChannelID: cm.ChannelId,
		},
		UserID: cm.UserId,
	}
}
