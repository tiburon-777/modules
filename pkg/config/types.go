package config

type Calendar struct {
	GRPC    Server
	HTTP    Server
	API     Server
	Logger  Logger
	Storage Storage
}

type Scheduler struct {
	Rabbitmq Rabbit
	Storage  Storage
	Logger   Logger
}

type Sender struct {
	Rabbitmq Rabbit
	Logger   Logger
}

type Server struct {
	Address string
	Port    string
}

type Rabbit struct {
	Login    string
	Pass     string
	Address  string
	Port     string
	Exchange string
	Queue    string
	Key      string
}

type Logger struct {
	File       string
	Level      string
	MuteStdout bool
}

type Storage struct {
	InMemory bool
	SQLHost  string
	SQLPort  string
	SQLDbase string
	SQLUser  string
	SQLPass  string
}
