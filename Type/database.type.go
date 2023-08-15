package Type

type Config struct {
	Host               string
	Port               string
	User               string
	Password           string
	Database           string
	SetMaxIdleConns    int
	SetMaxOpenConns    int
	SetConnMaxLifetime int
}
