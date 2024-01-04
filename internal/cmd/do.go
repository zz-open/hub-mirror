package cmd

import (
	"context"
	"log"
	"os"

	"github.com/zz-open/hub-mirror/internal/config"
	"github.com/zz-open/hub-mirror/internal/docker"
)

func Do() {
	log.Printf("mirrors: %+v, platform: %+v\n", config.C.Content, config.C.HubMirrors.Platform)

	log.Println("初始化 Docker 客户端")
	cli, err := docker.NewCli(context.Background(), config.C.Repository, config.C.Username, config.C.Password, os.Stdout)
	if err != nil {
		panic(err)
	}

}
