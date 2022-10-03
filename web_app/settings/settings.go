package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	viper.AddConfigPath(".")             // 还可以在工作目录中查找配置
	err = viper.ReadInConfig()           // 查找并读取配置文件
	if err != nil {                      // 处理读取配置文件的错误
		//读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed:%s", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了....")
	})
	return err
}
