package appconfig

type Env struct {
	DB_USER              string
	DB_PASSWORD          string
	DB_HOST              string
	DB_NAME              string
	DB_PORT              string
	SECRET_KEY           string
	ACCESS_SECRET        string
	REFRESH_SECRET       string
	TIMEOUT              string
	ACCESS_TOK_EXP       int
	REFRESH_TOK_EXP      int
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	DATABASE_URL         string
}
