/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/12 17:25
 */

package social

import (
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/social/rpc/socialclient"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
)

var (
	client zrpc.Client
	logic  social.SocialClient
	dsn    string
)

type MySQL struct {
	Database string `mapstructure:"Database"`
	Account  string `mapstructure:"Account"`
	Password string `mapstructure:"Password"`
	Host     string `mapstructure:"Host"`
	Port     int16  `mapstructure:"Port"`
	Options  string `mapstructure:"Options"`
}

func init() {
	client = zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:8004",
	})
	logic = socialclient.NewSocial(client)

	//读取测试配置文件
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("../social/config")

	err := v.ReadInConfig() // 读取配置文件
	if err != nil {
		// 处理读取配置文件失败的情况
		log.Panicln("读取测试配置文件失败，请重试", err)
	}

	//绑定成为结构体
	var mysql MySQL
	err = v.UnmarshalKey("MySQL", &mysql)
	if err != nil {
		log.Panicln("读取测试数据库失败，请重试", err)
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", mysql.Account, mysql.Password, mysql.Host, mysql.Port, mysql.Database, mysql.Options)
	log.Println(dsn)
}

// GetTestDB 获取测试数据库的连接
func GetTestDB() sqlx.SqlConn {
	conn := sqlx.NewMysql(dsn)
	return conn
}
