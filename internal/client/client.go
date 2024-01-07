package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/zz-open/hub-mirror/internal/config"
)

type HubMirrorClient struct {
	client   *client.Client
	server   string
	username string
	password string
	auth     string
	writer   io.Writer
}

type TransferInfo struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// 可选参数 ==============================

type Option func(*HubMirrorClient)

func WithServer(server string) Option {
	return func(cli *HubMirrorClient) {
		cli.server = server
	}
}

func WithUsername(username string) Option {
	return func(cli *HubMirrorClient) {
		cli.username = username
	}
}

func WithPassword(password string) Option {
	return func(cli *HubMirrorClient) {
		cli.password = password
	}
}

func WithWriter(writer io.Writer) Option {
	return func(cli *HubMirrorClient) {
		cli.writer = writer
	}
}

// 构造函数 ==============================

func NewHubMirrorClient(ctx context.Context, opts ...Option) (*HubMirrorClient, error) {
	c := &HubMirrorClient{}
	for _, o := range opts {
		o(c)
	}

	if c.username == "" || c.password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	if c.writer == nil {
		c.writer = os.Stdout
	}

	authConfig := c.newAuthConfig()
	auth, err := registry.EncodeAuthConfig(*authConfig)
	if err != nil {
		return nil, err
	}

	c.auth = auth

	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	_, err = dockerClient.RegistryLogin(ctx, *authConfig)
	if err != nil {
		return nil, err
	}

	c.client = dockerClient
	return c, nil
}

func (c *HubMirrorClient) newAuthConfig() *registry.AuthConfig {
	authConfig := &registry.AuthConfig{
		Username:      c.username,
		Password:      c.password,
		ServerAddress: c.server,
	}

	return authConfig
}

// 业务函数 ==============================

type DockerError struct {
	ErrorDetail string `json:"errorDetail"`
	Error       string `json:"error"`
}

func (c *HubMirrorClient) read(ctx context.Context, rd io.Reader) error {
	if rd == nil {
		return errors.New("rd 不能为Nil")
	}

	// docker pull 的过程中会输出json格式的信息，解析这些输出判断docker执行过程有没有出错
	// 例如：{"errorDetail":{"message":"denied: requested access to the resource is denied"},"error":"denied: requested access to the resource is denied"}
	var dockerError DockerError
	bufioReader := bufio.NewReader(rd)
	for {
		n, err := bufioReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if c.writer != nil {
			_, _ = c.writer.Write(n)
		}

		_ = json.Unmarshal(n, &dockerError)
		if dockerError.Error != "" {
			return errors.New(dockerError.Error)
		}

		if dockerError.ErrorDetail != "" {
			return errors.New(dockerError.ErrorDetail)
		}
	}

	return nil
}

func (c *HubMirrorClient) PullImage(ctx context.Context, image, platform string) error {
	reader, err := c.client.ImagePull(ctx, image, types.ImagePullOptions{Platform: platform})
	if err != nil {
		return err
	}

	defer reader.Close()

	return c.read(ctx, reader)
}

func (c *HubMirrorClient) PushImage(ctx context.Context, image, platform string) error {
	reader, err := c.client.ImagePush(ctx, image, types.ImagePushOptions{
		RegistryAuth: c.auth,
		Platform:     platform,
	})

	if err != nil {
		return err
	}

	defer reader.Close()

	return c.read(ctx, reader)
}

func (c *HubMirrorClient) TagImage(ctx context.Context, source, target string) error {
	return c.client.ImageTag(ctx, source, target)
}

func (c *HubMirrorClient) TransferImage(ctx context.Context, mirror, platform string) (*TransferInfo, error) {
	transferInfo, err := c.transfer(mirror)
	if err != nil {
		return nil, err
	}

	err = c.PullImage(ctx, transferInfo.Source, platform)
	if err != nil {
		return nil, err
	}

	err = c.TagImage(ctx, transferInfo.Source, transferInfo.Target)
	if err != nil {
		return nil, err
	}

	err = c.PushImage(ctx, transferInfo.Target, platform)
	if err != nil {
		return nil, err
	}

	return transferInfo, nil
}

func (c *HubMirrorClient) transfer(source string) (*TransferInfo, error) {
	if source == "" {
		return nil, errors.New("source 不能为空")
	}

	target := source
	splits := strings.Split(source, config.MIRROR_SEPERATOR)
	if len(splits) > 1 {
		source = splits[0]
		target = splits[1]
	}

	var prefix string
	// hub.docker.com
	if c.server == "" {
		prefix = c.username
	} else {
		// 其他
		prefix = c.server
	}

	target = fmt.Sprint(prefix, "/", target)
	return &TransferInfo{
		Source: source,
		Target: target,
	}, nil
}
