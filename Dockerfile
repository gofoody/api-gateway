FROM alpine:3.10

RUN mkdir /app
WORKDIR /app
ADD ./.out/api-gateway /app/api-gateway

CMD ["./api-gateway"]
