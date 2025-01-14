// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package appservices

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/incoming"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

var ErrNotABot = errors.New("not a bot")
var ErrIsABot = errors.New("is a bot")

type Service interface {
	// Subscriptions

	Subscribe(r *incoming.Request, sub apps.Subscription) error
	GetSubscriptions(r *incoming.Request, appID apps.AppID, actingUserID string) ([]apps.Subscription, error)
	Unsubscribe(r *incoming.Request, sub apps.Subscription) error

	// KV

	KVSet(r *incoming.Request, appID apps.AppID, actingUserID, prefix, id string, data []byte) (bool, error)
	KVGet(r *incoming.Request, appID apps.AppID, actingUserID, prefix, id string) ([]byte, error)
	KVDelete(r *incoming.Request, appID apps.AppID, actingUserID, prefix, id string) error
	KVList(r *incoming.Request, appID apps.AppID, actingUserID, namespace string, processf func(key string) error) error

	// Remote (3rd party) OAuth2

	StoreOAuth2App(r *incoming.Request, appID apps.AppID, actingUserID string, data []byte) error
	StoreOAuth2User(r *incoming.Request, AppID apps.AppID, actingUserID string, data []byte) error
	GetOAuth2User(r *incoming.Request, appID apps.AppID, actingUserID string) ([]byte, error)
}

type AppServices struct {
	conf  config.Service
	store *store.Service
}

var _ Service = (*AppServices)(nil)

func NewService(conf config.Service, store *store.Service) *AppServices {
	return &AppServices{
		conf:  conf,
		store: store,
	}
}

func (a *AppServices) ensureFromUser(mattermostUserID string) error {
	if mattermostUserID == "" {
		return utils.NewUnauthorizedError("not logged in")
	}
	mmuser, err := a.conf.MattermostAPI().User.Get(mattermostUserID)
	if err != nil {
		return err
	}
	if mmuser.IsBot {
		return errors.Wrap(ErrIsABot, mmuser.GetDisplayName(model.ShowNicknameFullName))
	}
	return nil
}
