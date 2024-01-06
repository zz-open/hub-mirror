package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zz-open/hub-mirror/internal/config"
)

func NewTestHubMirrorClientWithServer() (*HubMirrorClient, error) {
	username := "xxx"
	password := "xxx"
	server := "registry.cn-hangzhou.aliyuncs.com/zzimage"

	ctx := context.Background()
	cli, err := NewHubMirrorClient(
		ctx,
		WithUsername(username),
		WithPassword(password),
		WithServer(server),
	)

	return cli, err
}

// go test -v -run TestNewHubMirrorClientWithServer
func TestNewHubMirrorClientWithServer(t *testing.T) {
	cli, err := NewTestHubMirrorClientWithServer()

	assert.NotNil(t, cli)
	assert.Nil(t, err)
}

// go test -v -run TestTransferWithServer
func TestTransferWithServer(t *testing.T) {
	cli, err := NewTestHubMirrorClientWithServer()

	assert.NotNil(t, cli)
	assert.Nil(t, err)

	prefix := cli.server + "/"

	output, err := cli.transfer("")
	assert.Nil(t, output)

	source := "registry.k8s.io/kube-apiserver"
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, source, output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "registry.k8s.io/kube-apiserver"), output.Target)

	source = "registry.k8s.io/kube-apiserver:v1.27.4"
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, source, output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "registry.k8s.io/kube-apiserver:v1.27.4"), output.Target)

	source = fmt.Sprint("registry.k8s.io/kube-apiserver:v1.27.4", config.MIRROR_SEPERATOR, "aaa/bbb")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "registry.k8s.io/kube-apiserver:v1.27.4", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "aaa/bbb"), output.Target)

	source = fmt.Sprint("registry.k8s.io/kube-apiserver:v1.27.4", config.MIRROR_SEPERATOR, "aaa/bbb:v1.0.0")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "registry.k8s.io/kube-apiserver:v1.27.4", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "aaa/bbb:v1.0.0"), output.Target)

	source = fmt.Sprint("nginx@sha256:123456", config.MIRROR_SEPERATOR, "nginx")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "nginx@sha256:123456", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "nginx"), output.Target)

	source = fmt.Sprint("nginx@sha256:123456", config.MIRROR_SEPERATOR, "nginx:mytag")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "nginx@sha256:123456", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "nginx:mytag"), output.Target)
}

func NewTestHubMirrorClientWithoutServer() (*HubMirrorClient, error) {
	username := "xxx"
	password := "xxx"

	ctx := context.Background()
	cli, err := NewHubMirrorClient(
		ctx,
		WithUsername(username),
		WithPassword(password),
	)

	return cli, err
}

// go test -v -run TestNewHubMirrorClientWithoutServer
func TestNewHubMirrorClientWithoutServer(t *testing.T) {
	cli, err := NewTestHubMirrorClientWithoutServer()

	assert.NotNil(t, cli)
	assert.Nil(t, err)
}

// go test -v -run TestTransferWithoutServer
func TestTransferWithoutServer(t *testing.T) {
	cli, err := NewTestHubMirrorClientWithoutServer()

	assert.NotNil(t, cli)
	assert.Nil(t, err)

	prefix := cli.username + "/"

	output, err := cli.transfer("")
	assert.Nil(t, output)

	source := "registry.k8s.io/kube-apiserver"
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, source, output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "registry.k8s.io/kube-apiserver"), output.Target)

	source = "registry.k8s.io/kube-apiserver:v1.27.4"
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, source, output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "registry.k8s.io/kube-apiserver:v1.27.4"), output.Target)

	source = fmt.Sprint("registry.k8s.io/kube-apiserver:v1.27.4", config.MIRROR_SEPERATOR, "aaa/bbb")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "registry.k8s.io/kube-apiserver:v1.27.4", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "aaa/bbb"), output.Target)

	source = fmt.Sprint("registry.k8s.io/kube-apiserver:v1.27.4", config.MIRROR_SEPERATOR, "aaa/bbb:v1.0.0")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "registry.k8s.io/kube-apiserver:v1.27.4", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "aaa/bbb:v1.0.0"), output.Target)

	source = fmt.Sprint("nginx@sha256:123456", config.MIRROR_SEPERATOR, "nginx")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "nginx@sha256:123456", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "nginx"), output.Target)

	source = fmt.Sprint("nginx@sha256:123456", config.MIRROR_SEPERATOR, "nginx:mytag")
	output, err = cli.transfer(source)
	assert.Nil(t, err)
	assert.Equal(t, "nginx@sha256:123456", output.Source)
	assert.Equal(t, fmt.Sprint(prefix, "nginx:mytag"), output.Target)
}
