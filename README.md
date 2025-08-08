# go-opentelemetry-poc

POC de instrumentação Go com OpenTelemetry e [go-wslogger](https://github.com/thiagozs/go-wslogger).

## Variáveis de ambiente com direnv

Para facilitar o gerenciamento das variáveis de ambiente, utilize o [direnv](https://direnv.net/):

1. Instale o direnv:

    ```sh
    sudo apt install direnv
    # ou veja instruções para seu sistema em https://direnv.net/docs/installation.html
    ```

2. Copie o exemplo:

   ```sh
   cp .envrc.example .envrc
   direnv allow
   ```

3. As variáveis serão carregadas automaticamente ao entrar na pasta do projeto.

## Como rodar

```sh
docker-compose up -d
# Em outro terminal:
go run main.go
```

Acesse [http://localhost:8080/work](http://localhost:8080/work) para gerar traces.

## Observabilidade

- Exportação OTLP para o collector (gRPC ou HTTP)
- Visualização no Jaeger [http://localhost:16686](http://localhost:16686)

## Requisitos

- Go 1.23+
- Docker Compose

## Autor

- desenvolvido por [Thiago Zilli Sarmento](https://github.com/thiagozs) ❤️
