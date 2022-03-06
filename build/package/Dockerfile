FROM golang:1.17 AS api-builder
LABEL stage=intermediate
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM node:17 AS web-builder
LABEL stage=intermediate
COPY web /web
WORKDIR /
RUN npx tailwindcss -c web/tailwind.config.js -i web/src/input.css -o web/src/css/styles.css

FROM alpine:latest
LABEL maintainer="Hendrik Jonas Schlehlein <hendrik.schlehlein@gmail.com>"
LABEL org.opencontainers.image.source https://github.com/julis-bw-nw/julis-service-register-app
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=api-builder /main /
COPY --from=web-builder /web /web
RUN chmod +x ./main
ENTRYPOINT [ "./main", "run" ]