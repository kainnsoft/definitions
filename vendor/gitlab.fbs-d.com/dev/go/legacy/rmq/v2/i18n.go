package rmq

var messages = `
{
	"error.no-connection": "Not connected",
	"error.invalid-param": "Invalid method parameter",
	"error.connection": "Connection error",
	"error.create-exchange": "Create exchange error",
	"error.create-queue": "Create queue error",
	"error.create-channel": "Create channel error",
	"error.bind-queue": "Bind queue error",
	"error.unbind-queue": "Unbind queue error",
	"error.consume": "Consume error",
	"error.delete-queue": "Delete queue error",
	"error.publish": "Publish error",
	"error.trigger-event": "Trigger event error",
	"error.request-timeout": "Request timeout",
	"error.internal": "Internal error",
	"error.empty-queue-name": "queue name is empty",
	"error.empty-uuid": "uuid is empty",
	"error.empty-body": "body is empty",
	"error.empty-exchange-name": "exchange name is empty"
}
`
