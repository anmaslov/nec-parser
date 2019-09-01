package config

// Конфигурация
type Config struct {
	// LogLevel Уровень логирования
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn, info, debug" env:"NP_LOG_LEVEL" default:"info"`
	// LogJSON флаг устанавливающий JSON-формат логов
	LogJSON bool `long:"log-json" description:"Enable force log format JSON" env:"NP_LOG_JSON"`
	// DbServer Сервер БД
	DbAddress string `long:"db-address" description:"Database address Server:port" env:"NP_DB_ADDRESS" required:"true"`
	// DbName Имя базы данных
	DbName string `long:"db-name" description:"Database Name" env:"NP_DB_NAME" required:"true"`
	// DbUser Пользователь для подключения к БД
	DbUser string `long:"db-user" description:"Database User" env:"NP_DB_USER" required:"true"`
	// DbPassword Пароль для соединения с БД
	DbPassword string `long:"db-password" description:"Database Password" env:"NP_DB_PASSWORD"`
	// DbTimeOut таймаут соединения с БД
	DbTimeOut int `long:"db-time-out" description:"Database time out" env:"NP_DB_TIMEOUT" default:"60"`
}

// Db конфигурация базы данных
type Db struct {
	Host     string
	Dbname   string
	Username string
	Password string
}
