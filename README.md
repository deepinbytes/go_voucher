

[![Build Status](https://travis-ci.com/deepinbytes/go_voucher.svg?branch=master)](https://travis-ci.com/deepinbytes/go_voucher)
[![Coverage Status](https://coveralls.io/repos/github/deepinbytes/go_voucher/badge.svg?branch=master)](https://coveralls.io/github/deepinbytes/go_voucher?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/deepinbytes/go_voucher)](https://goreportcard.com/report/github.com/deepinbytes/go_voucher)



# Go REST Voucher Pool

A voucher pool is a collection of voucher codes that can be used by customers to get discounts on website. Each code may only be used once, and we would like to know when it was used by the customer. 
this internal application exposes REST endpoints to manage the voucher pool, offers and users.


**Used libraries:**
- [gin](https://github.com/gin-gonic)
- [gin-swagger](https://github.com/swaggo/gin-swagger)
- [gorm](https://gorm.io/docs/)
- [godotenv](https://pkg.go.dev/github.com/joho/godotenv?tab=doc)
- [testify](https://github.com/stretchr/testify)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)

---

Swagger Doc at http://localhost:3000/swagger/index.html

### Run locally

Create `.env` at root, i.e.
```sh
DB_HOST=pg << localhost in case of running locally
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=local-dev-db

ENV=development

APP_PORT=3000
APP_HOST=http://localhost
```

Run
```sh
# Terminal 1
docker-compose up        # docker-compose up (Run postgres and application)
docker-compose down      # docker-compose down (Shutdown)

```

Test
```
go test -v -cover ./...               # Run go test
```

---

### Todo

- [ ] Access Control
- [ ] Input Validations
- [ ] Custom Error messages
- [ ] Logger
- [ ] More unit tests


---

