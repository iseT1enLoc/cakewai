package appconfig

type Env struct {
	DB_USER        string
	DB_PASSWORD    string
	DB_HOST        string
	DB_NAME        string
	DB_PORT        string
	SECRET_KEY     string
	ACCESS_SECRET  string
	REFRESH_SECRET string
	TIMEOUT        string
}
