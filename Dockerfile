FROM golang:alpine3.14 as build-env
 
# Set environment variable
ENV APP_NAME apex-discord-bot
ENV CMD_PATH main.go
 
# Copy application data into image
COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME
 
# Build application
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
 
# Run Stage
FROM alpine:3.14
 
# Set environment variable
ENV APP_NAME apex-discord-bot
 
# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
COPY .env .
 
ENTRYPOINT ["./apex-discord-bot"]


