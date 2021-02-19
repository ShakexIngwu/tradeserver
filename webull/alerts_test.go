package webull

import (
	"os"
	"testing"

	model "github.com/ShakexIngwu/tradeserver/webullmodel"
	"github.com/stretchr/testify/assert"
)

func TestGetAlerts(t *testing.T) {
	if os.Getenv("WEBULL_USERNAME") == "" {
		t.Skip("No username set")
		return
	}
	asrt := assert.New(t)
	c, err := NewClient(&Credentials{
		Username:    os.Getenv("WEBULL_USERNAME"),
		Password:    os.Getenv("WEBULL_PASSWORD"),
		AccountType: model.AccountType(2),
		DeviceName:  deviceName(),
	})
	asrt.Empty(err)
	asrt.NotNil(c)
	alerts, err := c.GetAlerts()
	asrt.NotEmpty(alerts)
	asrt.Empty(err)
}
