package mysql

type connect struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}
type Config struct {
	Master connect `yaml:"master"`
	Slaves []connect `yaml:"slaves"` // 从库配置
	MaxOpenConns    int `yaml:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"` // 连接池最大空闲连接数
	LogLevel        string `yaml:"log_level"`
}