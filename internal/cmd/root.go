package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"github.com/zz-open/hub-mirror/internal/config"
)

var (
	Cmd *cobra.Command

	cfgFile    string
	content    string
	maximum    int
	repository string
	username   string
	password   string
	output     string
)

func init() {
	cobra.OnInitialize(initConfig)
	Cmd = &cobra.Command{
		Use:   "hub-mirror",
		Short: "pull镜像上传到私有云仓库",
		Long:  `pull镜像上传到私有云仓库`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("config.C, %#v", config.C)
			Do()
		},
	}

	// 禁用自动补全子命令
	Cmd.CompletionOptions.DisableDefaultCmd = true
	Cmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "", "配置文件")
	Cmd.PersistentFlags().StringVarP(&content, "content", "", "", "原始镜像，格式为：{ \"platform\": \"\", \"hub-mirror\": [] }")
	Cmd.PersistentFlags().IntVarP(&maximum, "maximum", "", 20, "镜像数量上限")
	Cmd.PersistentFlags().StringVarP(&repository, "repo", "", "hub.docker.com", "推送仓库地址，默认为 hub.docker.com")
	Cmd.PersistentFlags().StringVarP(&username, "username", "", "", "仓库用户名")
	Cmd.PersistentFlags().StringVarP(&password, "password", "", "", "仓库密码")
	Cmd.PersistentFlags().StringVarP(&output, "output", "", "output.sh", "结果输出路径")

	// 绑定viper
	_ = viper.BindPFlag("Content", Cmd.PersistentFlags().Lookup("content"))
	_ = viper.BindPFlag("Maximum", Cmd.PersistentFlags().Lookup("maximum"))
	_ = viper.BindPFlag("Repository", Cmd.PersistentFlags().Lookup("repo"))
	_ = viper.BindPFlag("Username", Cmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("Password", Cmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("Output", Cmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	// 优先配置文件方式启动，其次命令行参数
	if cfgFile == "" {
		log.Println("采用命令行方式启动")
	} else {
		log.Println("采用配置文件方式启动")
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
		log.Fatalln("unable to decode into config.C", err)
	}

	if config.C.Content == "" {
		log.Fatalln("镜像列表不能为空")
	}

	var hubMirrors config.HubMirrors
	err = json.Unmarshal([]byte(config.C.Content), &hubMirrors)
	if err != nil {
		panic(err)
	}

	if len(hubMirrors.Content) > maximum {
		log.Fatalln("镜像数量超出最大限制")
	}

	config.C.HubMirrors = &hubMirrors
}

func Execute() {
	if err := Cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
