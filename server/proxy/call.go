// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/upstream"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

// CallResponse contains everything the CallResponse struct contains, plus some additional
// data for the client, such as information about the App's bot account.
//
// Apps will use the CallResponse struct to respond to a CallRequest, and the proxy will
// decorate the response using the CallResponse to provide additional information.
type CallResponse struct {
	apps.CallResponse

	// Used to provide info about the App to client, e.g. the bot user id
	AppMetadata AppMetadataForClient `json:"app_metadata"`
}

type AppMetadataForClient struct {
	BotUserID   string `json:"bot_user_id,omitempty"`
	BotUsername string `json:"bot_username,omitempty"`
}

func NewProxyCallResponse(response apps.CallResponse) CallResponse {
	return CallResponse{
		CallResponse: response,
	}
}

func (r CallResponse) WithMetadata(metadata AppMetadataForClient) CallResponse {
	r.AppMetadata = metadata
	return r
}

func (p *Proxy) Call(in Incoming, creq apps.CallRequest) CallResponse {
	if creq.Context.AppID == "" {
		return NewProxyCallResponse(apps.NewErrorResponse(
			utils.NewInvalidError("app_id is not set in Context, don't know what app to call")))
	}

	app, err := p.store.App.Get(creq.Context.AppID)
	if err != nil {
		return NewProxyCallResponse(apps.NewErrorResponse(err))
	}

	cresp, _ := p.callApp(in, *app, creq)
	return NewProxyCallResponse(cresp).WithMetadata(AppMetadataForClient{
		BotUserID:   app.BotUserID,
		BotUsername: app.BotUsername,
	})
}

func (p *Proxy) call(in Incoming, app apps.App, call apps.Call, cc *apps.Context, valuePairs ...interface{}) apps.CallResponse {
	values := map[string]interface{}{}
	for len(valuePairs) > 0 {
		if len(valuePairs) == 1 {
			return apps.NewErrorResponse(
				errors.Errorf("mismatched parameter count, no value for %v", valuePairs[0]))
		}
		key, ok := valuePairs[0].(string)
		if !ok {
			return apps.NewErrorResponse(
				errors.Errorf("mismatched type %T for key %v, expected string", valuePairs[0], valuePairs[0]))
		}
		values[key] = valuePairs[1]
		valuePairs = valuePairs[2:]
	}

	if cc == nil {
		cc = &apps.Context{}
	}
	cresp, _ := p.callApp(in, app, apps.CallRequest{
		Call:    call,
		Context: *cc,
		Values:  values,
	})
	return cresp
}

func (p *Proxy) callApp(in Incoming, app apps.App, creq apps.CallRequest) (apps.CallResponse, error) {
	respondErr := func(err error) (apps.CallResponse, error) {
		return apps.NewErrorResponse(err), err
	}

	conf, _, log := p.conf.Basic()
	log = log.With("app_id", app.AppID)

	if !p.appIsEnabled(app) {
		return respondErr(errors.Errorf("%s is disabled", app.AppID))
	}

	if creq.Path[0] != '/' {
		return respondErr(utils.NewInvalidError("call path must start with a %q: %q", "/", creq.Path))
	}
	cleanPath, err := utils.CleanPath(creq.Path)
	if err != nil {
		return respondErr(err)
	}
	creq.Path = cleanPath

	up, err := p.upstreamForApp(app)
	if err != nil {
		return respondErr(err)
	}

	cc := creq.Context
	cc = in.updateContext(cc)
	creq.Context, err = p.expandContext(in, app, &cc, creq.Expand)
	if err != nil {
		return respondErr(err)
	}

	cresp, err := upstream.Call(up, app, creq)
	if err != nil {
		return cresp, err
	}
	if cresp.Type == "" {
		cresp.Type = apps.CallResponseTypeOK
	}

	if cresp.Form != nil {
		if cresp.Form.Icon != "" {
			icon, err := normalizeStaticPath(conf, cc.AppID, cresp.Form.Icon)
			if err != nil {
				log.WithError(err).Debugw("Invalid icon path in form. Ignoring it.", "icon", cresp.Form.Icon)
				cresp.Form.Icon = ""
			} else {
				cresp.Form.Icon = icon
			}
			clean, problems := cleanForm(*cresp.Form)
			for _, prob := range problems {
				log.WithError(prob).Debugw("invalid form")
			}
			cresp.Form = &clean
		}
	}

	return cresp, nil
}

// normalizeStaticPath converts a given URL to a absolute one pointing to a static asset if needed.
// If icon is an absolute URL, it's not changed.
// Otherwise assume it's a path to a static asset and the static path URL prepended.
func normalizeStaticPath(conf config.Config, appID apps.AppID, icon string) (string, error) {
	if !strings.HasPrefix(icon, "http://") && !strings.HasPrefix(icon, "https://") {
		cleanIcon, err := utils.CleanStaticPath(icon)
		if err != nil {
			return "", errors.Wrap(err, "invalid icon path")
		}

		icon = conf.StaticURL(appID, cleanIcon)
	}

	return icon, nil
}

func (p *Proxy) GetStatic(appID apps.AppID, path string) (io.ReadCloser, int, error) {
	app, err := p.store.App.Get(appID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, utils.ErrNotFound) {
			status = http.StatusNotFound
		}
		return nil, status, err
	}

	return p.getStatic(*app, path)
}

func (p *Proxy) getStatic(app apps.App, path string) (io.ReadCloser, int, error) {
	up, err := p.upstreamForApp(app)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return up.GetStatic(app, path)
}