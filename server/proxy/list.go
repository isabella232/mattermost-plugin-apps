// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"sort"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/incoming"
)

func (p *Proxy) GetManifest(r *incoming.Request, appID apps.AppID) (*apps.Manifest, error) {
	return p.store.Manifest.Get(r, appID)
}

func (p *Proxy) GetInstalledApp(r *incoming.Request, appID apps.AppID) (*apps.App, error) {
	return p.store.App.Get(r, appID)
}

func (p *Proxy) GetInstalledApps(r *incoming.Request) []apps.App {
	installed := p.store.App.AsMap(r)
	out := []apps.App{}
	for _, app := range installed {
		out = append(out, app)
	}

	// Sort result alphabetically, by display name.
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].DisplayName) < strings.ToLower(out[j].DisplayName)
	})

	return out
}

func (p *Proxy) GetListedApps(r *incoming.Request, filter string, includePluginApps bool) []apps.ListedApp {
	conf := p.conf.Get()
	out := []apps.ListedApp{}

	for _, m := range p.store.Manifest.AsMap(r) {
		if !appMatchesFilter(m, filter) {
			continue
		}

		if !includePluginApps && m.Contains(apps.DeployPlugin) {
			continue
		}

		marketApp := apps.ListedApp{
			Manifest: m,
		}

		if m.Icon != "" {
			marketApp.IconURL = conf.StaticURL(m.AppID, m.Icon)
		}

		app, _ := p.store.App.Get(r, m.AppID)
		if app != nil {
			marketApp.Installed = true
			marketApp.Enabled = !app.Disabled
			marketApp.Labels = []model.MarketplaceLabel{{
				Name:        "Experimental",
				Description: "Apps are marked as experimental and not meant for production use. Please use with caution.",
				URL:         "",
			}}

			if !marketApp.Enabled {
				marketApp.Labels = append(marketApp.Labels, model.MarketplaceLabel{
					Name:        "Disabled",
					Description: "This app is disabled.",
					URL:         "",
				})
			}
		}
		out = append(out, marketApp)
	}

	// Sort result alphabetically, by display name.
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].Manifest.DisplayName) < strings.ToLower(out[j].Manifest.DisplayName)
	})

	return out
}

// Copied from Mattermost Server
func appMatchesFilter(manifest apps.Manifest, filter string) bool {
	filter = strings.TrimSpace(strings.ToLower(filter))

	if filter == "" {
		return true
	}

	if strings.ToLower(string(manifest.AppID)) == filter {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.DisplayName), filter) {
		return true
	}

	if strings.Contains(strings.ToLower(manifest.Description), filter) {
		return true
	}

	return false
}
