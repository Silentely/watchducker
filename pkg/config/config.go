package config

import (
	"fmt"

	"watchducker/pkg/logger"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	useLabel       bool     `mapstructure:"label"`
	useNoRestart   bool     `mapstructure:"no-restart"`
	cronExpression string   `mapstructure:"cron"`
	containerNames []string `mapstructure:"-"` // 位置参数，不通过mapstructure绑定
	logLevel       string   `mapstructure:"log_level"`
}

// 全局配置实例（只读，初始化后不可修改）
var globalConfig *Config

// Load 加载配置并初始化全局配置实例
func Load() error {
	// 如果已经初始化，直接返回
	if globalConfig != nil {
		return nil
	}

	config, err := loadConfig()
	if err != nil {
		return err
	}

	globalConfig = config
	return nil
}

// Get 获取全局配置实例（只读访问）
func Get() *Config {
	return globalConfig
}

// UseLabel 获取 UseLabel 配置
func (c *Config) UseLabel() bool {
	return c.useLabel
}

// UseNoRestart 获取 UseNoRestart 配置
func (c *Config) UseNoRestart() bool {
	return c.useNoRestart
}

// CronExpression 获取 CronExpression 配置
func (c *Config) CronExpression() string {
	return c.cronExpression
}

// ContainerNames 获取 ContainerNames 配置
func (c *Config) ContainerNames() []string {
	return c.containerNames
}

// LogLevel 获取 LogLevel 配置
func (c *Config) LogLevel() string {
	return c.logLevel
}

// loadConfig 执行实际的配置加载逻辑
func loadConfig() (*Config, error) {
	// 创建 Viper 实例
	v := viper.New()
	v.SetEnvPrefix("WATCHDUCKER")
	v.AutomaticEnv()

	// 设置 Viper 默认值
	v.SetDefault("label", false)
	v.SetDefault("no-restart", false)
	v.SetDefault("cron", nil)

	// 设置命令行参数
	pflag.Bool("label", false, "检查所有带有 watchducker.update=true 标签的容器")
	pflag.Bool("no-restart", false, "只更新镜像，不重启容器")
	pflag.String("cron", "", "定时执行，使用标准 cron 表达式格式")

	// 解析命令行参数
	pflag.Parse()

	// 绑定命令行参数到 Viper
	v.BindPFlags(pflag.CommandLine)

	// 创建配置实例
	config := &Config{
		useLabel:       v.GetBool("label"),
		useNoRestart:   v.GetBool("no-restart"),
		cronExpression: v.GetString("cron"),
		// 获取位置参数（容器名称）
		containerNames: pflag.Args(),
		logLevel:       v.GetString("LOG_LEVEL"),
	}

	// 设置日志级别
	if config.logLevel != "" {
		logger.SetLevel(config.logLevel)
	}

	// 验证配置有效性
	if err := config.validate(); err != nil {
		PrintUsage()
		return nil, err
	}

	return config, nil
}

// Validate 验证配置的有效性
func (c *Config) validate() error {
	// 验证至少需要一种检查方式
	if len(c.containerNames) == 0 && !c.useLabel {
		return fmt.Errorf("必须指定容器名称或使用 --label 选项")
	}

	return nil
}

// PrintUsage 打印使用方法
func PrintUsage() {
	fmt.Println("\n使用方法:")
	fmt.Println("  watchducker [选项] [容器名称...]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  --label       检查所有带有 watchducker.update=true 标签的容器")
	fmt.Println("  --no-restart  只更新镜像，不重启容器")
	fmt.Println("  --cron       定时执行，使用标准 cron 表达式格式")
	fmt.Println()
	fmt.Println("环境变量:")
	fmt.Println("  WATCHDUCKER_USE_LABEL    等同于 --label 选项")
	fmt.Println("  WATCHDUCKER_NO_RESTART   等同于 --no-restart 选项")
	fmt.Println("  WATCHDUCKER_CRON         等同于 --cron 选项")
	fmt.Println("  WATCHDUCKER_LOG_LEVEL    设置日志级别 (DEBUG/INFO/WARN/ERROR)")
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  要检查的容器名称列表（支持多个）  <容器1> <容器2> ... ")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 检查指定容器")
	fmt.Println("  watchducker nginx redis mysql")
	fmt.Println()
	fmt.Println("  # 检查所有带有 watchducker.update=true 标签的容器")
	fmt.Println("  watchducker --label")
	fmt.Println()
	fmt.Println("  # 使用环境变量配置")
	fmt.Println("  export WATCHDUCKER_LOG_LEVEL=DEBUG")
	fmt.Println("  export WATCHDUCKER_USE_LABEL=true")
	fmt.Println("  export WATCHDUCKER_CRON=\"0 2 * * *\"")
	fmt.Println()
	fmt.Println("  # 定时执行示例")
	fmt.Println("  watchducker --cron \"0 2 * * *\" --label          # 每天凌晨2点检查所有标签容器")
	fmt.Println("  watchducker --cron \"*/30 * * * *\" nginx redis   # 每30分钟检查指定容器")
	fmt.Println("  watchducker --cron \"@daily\" --no-restart        # 每天执行，只检查不重启")
	fmt.Println()
	fmt.Println("说明:")
	fmt.Println("  - 如果指定了容器名称，则检查这些特定容器")
	fmt.Println("  - 如果使用了 --label 选项，则检查所有带有 watchducker.update=true 标签的容器")
	fmt.Println("  - 两者不能同时使用")
	fmt.Println("  - 使用 --cron 参数时，程序会持续运行并按计划执行")
	fmt.Println("  - 环境变量优先级低于命令行参数")
}
