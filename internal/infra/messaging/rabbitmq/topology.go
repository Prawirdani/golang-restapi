package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func SetupTopology(ch *amqp.Channel) error {
	if err := setupAuthTopology(ch); err != nil {
		return err
	}
	// Setup other topologies here ....

	return nil
}
