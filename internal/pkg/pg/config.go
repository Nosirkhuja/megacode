package pg

type PgConfig struct {
	name     string `envconfig:"NAME" required:"true"`
	password string `envconfig:"PASSWORD" required:"true"`
	port     string `envconfig:"PORT" required:"true"`
}
