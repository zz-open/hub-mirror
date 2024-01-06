package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zz-open/hub-mirror/internal/client"
	"github.com/zz-open/hub-mirror/internal/config"
)

var (
	Cmd *cobra.Command

	cfgFile    string
	server     string
	username   string
	password   string
	rawContent string
	maximum    int
	output     string
)

func init() {
	cobra.OnInitialize(initConfig)
	Cmd = &cobra.Command{
		Use:   "hub-mirror",
		Short: "加速国外镜像",
		Long:  `拉国外镜像上传到国内私有云仓库加速下载,支持registry.k8s.io, k8s.gcr.io, gcr.io, quay.io等`,
		Run: func(cmd *cobra.Command, args []string) {
			client.Do()
		},
	}

	// 禁用自动补全子命令
	Cmd.CompletionOptions.DisableDefaultCmd = true
	Cmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "", "配置文件")
	Cmd.PersistentFlags().StringVarP(&rawContent, "rawcontent", "r", "", "拉取镜像清单，格式为：[{\"platform\":\"\", \"mirrors\":[]}]")
	Cmd.PersistentFlags().StringVarP(&username, "username", "u", "", "推送仓库用户名")
	Cmd.PersistentFlags().StringVarP(&password, "password", "p", "", "推送仓库密码")
	Cmd.PersistentFlags().StringVarP(&server, "server", "", "server", "推送仓库地址,默认为 hub.docker.com")
	Cmd.PersistentFlags().IntVarP(&maximum, "maximum", "", config.MIRROR_MAMIMUM, "拉取镜像数量上限")
	Cmd.PersistentFlags().StringVarP(&output, "output", "o", "output.sh", "脚本输出路径")

	// cmd变量绑定到viper
	_ = viper.BindPFlag("Server", Cmd.PersistentFlags().Lookup("server"))
	_ = viper.BindPFlag("Username", Cmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("Password", Cmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("Maximum", Cmd.PersistentFlags().Lookup("maximum"))
	_ = viper.BindPFlag("Output", Cmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	var mode = config.MODE_FILE
	// 优先配置文件方式启动，其次命令行参数
	if cfgFile == "" {
		log.Println("采用命令行方式启动")
		mode = config.MODE_CMD_PARAMETER
	} else {
		log.Println("采用配置文件方式启动")
		mode = config.MODE_FILE
		_, err := os.Stat(cfgFile)
		if err != nil {
			log.Fatalln("配置文件不存在", err)
		}

		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("不能读取配置文件:", err)
		}
	}

	err := viper.Unmarshal(&config.C)
	if err != nil {
		log.Fatalln("解析到全局变量失败", err)
	}

	if config.C.Username == "" {
		log.Fatalln("推送仓库用户名不能为空")
	}

	if config.C.Password == "" {
		log.Fatalln("推送仓库密码不能为空")
	}

	mirrors := make([]config.Mirrors, 0)
	if mode == config.MODE_CMD_PARAMETER {
		if rawContent == "" {
			log.Fatalln("rawcontent必传")
		}

		err = json.Unmarshal([]byte(rawContent), &mirrors)
		if err != nil {
			log.Fatalln(err)
		}
	} else if mode == config.MODE_FILE {
		mirrors = config.C.Mirrors
	}

	if len(mirrors) == 0 {
		log.Fatalln("拉取镜像清单不能为空")
	}

	var imageNum int
	for _, v := range mirrors {
		for _, m := range v.Mirrors {
			if m == "" {
				continue
			}
			imageNum += 1
		}
	}

	if imageNum == 0 {
		log.Fatalln("拉取镜像清单不能为空")
	}

	if imageNum > maximum {
		log.Fatalln(fmt.Sprintf("最多拉取%d个镜像", maximum))
	}

	if mode == config.MODE_CMD_PARAMETER {
		config.C.Mirrors = mirrors
	}
}

func Execute() {
	if err := Cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
