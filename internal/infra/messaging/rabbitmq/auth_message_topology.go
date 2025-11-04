package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

var (
	AuthExchange = Exchange{
		Name:       "auth.direct",
		Type:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
	AuthEmailResetPasswordQueue = Queue{
		Name:       "auth.email.reset-password",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
)

const (
	ResetPasswordEmailRoutingKey = "email.reset-password"
)

// Setup auth exchange, queue and binds
func setupAuthTopology(ch *amqp.Channel) error {
	if err := AuthExchange.Declare(ch); err != nil {
		return err
	}

	if err := AuthEmailResetPasswordQueue.Declare(ch); err != nil {
		return err
	}

	return AuthExchange.BindQueue(
		ch,
		AuthEmailResetPasswordQueue,
		ResetPasswordEmailRoutingKey,
	)
}
