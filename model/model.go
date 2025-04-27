package model

type Configurations struct {
	Prefix string
	Server struct {
		Host string
		Port string
	}
	MysqlConf   MysqlConfig
	RedisConfig RedisConfig
}

type MysqlConfig struct {
	Username     string
	Password     string
	Net          string
	Address      string
	DatabaseName string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Employee struct {
	ID          string `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	CompanyName string `json:"company_name" db:"company_name"`
	Address     string `json:"address" db:"address"`
	City        string `json:"city" db:"city"`
	County      string `json:"county" db:"county"`
	Postal      string `json:"postal" db:"postal"`
	Phone       string `json:"phone" db:"phone"`
	Email       string `json:"email" db:"email"`
	Web         string `json:"web" db:"web"`
}

