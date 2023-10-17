# Chat API

This is a web chat API that allows users to create and manage chat rooms and send messages. The chat system also includes another API called Broadcaster, where users can subscribe to receive messages from a given chat room. The messages are sent from this API to Broadcaster through RabbitMQ.


## Endpoints

| Endpoint                     | Method | Protected | Description         |
|------------------------------| ------ |-----------|---------------------|
| `/api/v1/rooms`              | POST   | YES       | Create a room       |
| `/api/v1/rooms`              | GET    | YES       | Search rooms        |
| `/api/v1/rooms/{id}`         | GET    | YES       | Find a room by id   |
| `/api/v1/rooms/{id}`         | PUT    | YES       | Update a room       |
| `/api/v1/rooms/{id}`         | DELETE | YES       | Delete a room       |
| `/api/v1/rooms/{id}/send`    | POST   | YES       | Send a message      |
| `/api/v1/swagger/index.html` | GET    | NO        | API's documentation |
| `/api/v1/healthz`            | GET    | NO        | Health check        |

## Related repositories

- [Broadcaster API](https://github.com/sesaquecruz/go-chat-broadcaster)
- [Chat Infra](https://github.com/sesaquecruz/k8s-chat-infra)
- [Chat API Docker Hub](https://hub.docker.com/r/sesaquecruz/go-chat-api/tags)

## Tech Stack

- [Go](https://go.dev)
- [Gin](https://gin-gonic.com)
- [Postgres](https://www.postgresql.org)
- [RabbitMQ](https://www.rabbitmq.com)


## Contributing

Contributions are welcome! If you find a bug or would like to suggest an enhancement, please make a fork, create a new branch with the bugfix or feature, and submit a pull request.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) file for more information.
