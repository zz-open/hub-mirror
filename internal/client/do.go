package client

import (
	"context"
	"html/template"
	"log"
	"os"
	"sync"

	"github.com/zz-open/hub-mirror/internal/config"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

func Do() {
	log.Println("初始化 Docker 客户端")

	ctx := context.Background()
	hubMirrorClient, err := NewHubMirrorClient(
		ctx,
		WithUsername(config.C.Username),
		WithPassword(config.C.Password),
		WithServer(config.C.Server),
	)

	if err != nil {
		log.Fatalln(err)
	}

	transferInfos := make([]*TransferInfo, 0)
	for _, v := range config.C.Mirrors {
		v := v
		for _, source := range v.Mirrors {
			source := source
			if source == "" {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				defer mu.Unlock()

				// 转换
				output, err := hubMirrorClient.TransferImage(context.Background(), source, v.Platform)
				if err != nil {
					log.Println(source, "转换失败: ", err)
					return
				}

				transferInfos = append(transferInfos, output)
			}()
		}
	}

	wg.Wait()

	if len(transferInfos) == 0 {
		log.Fatalln("没有转换成功的镜像")
	}

	tmpl, err := template.New("output").ParseFiles("output.tpl")
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(config.C.Output)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()
	err = tmpl.Execute(f, map[string]interface{}{
		"Outputs": transferInfos,
		"Server":  config.C.Server,
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("SUCCESS")
}
