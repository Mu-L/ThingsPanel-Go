package app

import (
	"project/internal/query"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Application 结构体用于管理所有应用依赖
type Application struct {
	Config         *viper.Viper
	Logger         *logrus.Logger
	DB             *gorm.DB
	RedisClient    *redis.Client
	ServiceManager *ServiceManager
}

// NewApplication 创建并初始化应用
func NewApplication(options ...Option) (*Application, error) {
	app := &Application{
		Logger:         logrus.New(),
		ServiceManager: NewServiceManager(),
	}

	// 应用所有选项
	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	// 设置查询默认DB
	if app.DB != nil {
		query.SetDefault(app.DB)
	}

	return app, nil
}

// RegisterService 注册一个服务到应用程序
func (app *Application) RegisterService(service Service) {
	app.ServiceManager.RegisterService(service)
}

// Start 启动所有注册的服务
func (app *Application) Start() error {
	return app.ServiceManager.StartAll()
}

// Shutdown 优雅关闭所有资源
func (app *Application) Shutdown() {
	// 停止所有服务
	app.ServiceManager.StopAll()

	// 关闭数据库连接
	if app.RedisClient != nil {
		app.RedisClient.Close()
		app.Logger.Info("Redis connection closed")
	}

	// DB不需要显式关闭，gorm.DB没有Close方法

	app.Logger.Info("All resources have been successfully cleaned up")
}

// Wait 等待所有服务完成
func (app *Application) Wait() {
	app.ServiceManager.Wait()
}

// Option 定义应用程序初始化选项
type Option func(*Application) error
