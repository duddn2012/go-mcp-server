package testutils

type TestConfig struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	AllowedOrigins []string
}

// NewTestConfig는 테스트용 기본 설정을 반환합니다
func NewTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "postgres",
		DBPassword:     "postgres",
		DBName:         "test_db",
		ServerPort:     "8080",
		AllowedOrigins: []string{"http://localhost"},
	}
}
