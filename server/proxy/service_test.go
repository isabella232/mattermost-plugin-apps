package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/mocks/mock_store"
	"github.com/mattermost/mattermost-plugin-apps/server/mocks/mock_upstream"
	"github.com/mattermost/mattermost-plugin-apps/server/store"
	"github.com/mattermost/mattermost-plugin-apps/upstream"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

func newTestProxy(tb testing.TB, testApps []apps.App, ctrl *gomock.Controller) *Proxy {
	conf := config.NewTestConfigService(nil).WithMattermostConfig(model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: model.NewString("test.mattermost.com"),
		},
	})

	s, err := store.MakeService(utils.NewTestLogger(), conf, nil)
	require.NoError(tb, err)
	appStore := mock_store.NewMockAppStore(ctrl)
	s.App = appStore

	upstreams := map[apps.AppID]upstream.Upstream{}
	for i := range testApps {
		app := testApps[i]

		up := mock_upstream.NewMockUpstream(ctrl)

		// set up an empty OK call response
		b, _ := json.Marshal(apps.NewDataResponse(nil))
		reader := ioutil.NopCloser(bytes.NewReader(b))
		up.EXPECT().Roundtrip(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, nil)

		upstreams[app.Manifest.AppID] = up
		appStore.EXPECT().Get(gomock.Any(), app.AppID).Return(&app, nil)
	}

	p := &Proxy{
		store:            s,
		builtinUpstreams: upstreams,
		conf:             conf,
	}

	return p
}
