package config_test

import (
	"testing"

	"github.com/longbridgeapp/assert"
	"github.com/longportapp/openapi-go/config"
)

func Test_withConfigKey(t *testing.T) {
	var c, err = config.New(config.WithConfigKey("appKey", "appSecret", "accessToken"))
	assert.NoError(t, err)
	assert.Equal(t, "appKey", c.AppKey)
	assert.Equal(t, "appSecret", c.AppSecret)
	assert.Equal(t, "accessToken", c.AccessToken)
}
