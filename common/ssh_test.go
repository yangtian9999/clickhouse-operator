package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yangtian9999/clickhouse-operator/log"
	"github.com/yangtian9999/clickhouse-operator/model"
)

func TestRemoteExecute(t *testing.T) {
	log.InitLoggerConsole()
	cmd := "grep lo /proc/net/if_inet6 >/dev/null 2>&1;echo $?"
	sshOpts := SshOptions{
		User:             "root",
		Password:         "123456",
		Port:             22,
		Host:             "192.168.21.73",
		NeedSudo:         false,
		AuthenticateType: model.SshPasswordSave,
	}
	out, err := RemoteExecute(sshOpts, cmd)
	assert.Nil(t, err)
	assert.Equal(t, "0", out)
}
