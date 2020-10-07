package config

type Configuration struct {
	AMQPConnectionURL string
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
}

type Task struct {
	Number1 int
	Number2 int
}
