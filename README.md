# pintas - URL Shortener
#### built with [htmx](https://htmx.org/), [gin](https://github.com/gin-gonic/gin), [gorm](https://github.com/go-gorm/gorm) and [tailwindcss](https://tailwindcss.com/)

![web preview](https://i.imgur.com/AznmuXm.png)

## Environment Variables

| Name | Description | Example |
| --- | --- | --- |
| PINTAS_DB_DSN | PostgreSQL DSN string | host=localhost port=5432 user=postgres password=postgres dbname=pintas sslmode=disable |

## Development

```bash
tailwindcss -i frontend/styles/index.css -o frontend/static/styles.css -w
air
```

## Build

```bash
tailwindcss -i frontend/styles/index.css -o frontend/static/styles.css 
GIN_MODE=release go build -o pintas main.go
./pintas
```
