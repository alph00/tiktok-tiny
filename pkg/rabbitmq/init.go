package rabbitmq

import (
	"fmt"
	"log"

	"github.com/alph00/tiktok-tiny/pkg/viper"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	config = viper.Read("rabbitmq")
	// logger *zap.SugaredLogger
	conn  *amqp.Connection
	err   error
	MqUrl = fmt.Sprint("amqp://guest:guest@localhost:5672/")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
