**1. Migration**
Для простоты, сделал базу данных для всего приложения, а не каждого микросервиса.

goose -dir ./migrations postgres "user=postgres password=postgres dbname=value_exchanger sslmode=disable" up