# Go Boilerplate
Easily expandable ready to use golang boilerplate.
## Features
- Configuration by [viper](https://github.com/spf13/viper)
- API Routing by [chi](https://github.com/go-chi/chi)
- API Documentation [swagger](github.com/flowchartsman/swaggerui)
- Database usage [gorm](https://github.com/go-gorm/gorm)
- Database migrations [gormigrate](https://github.com/go-gormigrate/gormigrate)
- AMQP client [rabbitmq](https://github.com/rabbitmq/amqp091-go)
- Redis client [redis](https://github.com/redis/go-redis)
- Logging with [zerolog](https://github.com/rs/zerolog)
## How to use
Copy .env.example file and fill environment variables
```
cp .env.example .env
```
Run
```
make tidy
make run
```