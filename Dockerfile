FROM golang:alpine as build-env

WORKDIR /app

COPY . ./

RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /otpapp ./cmd/otpapp/main.go

FROM alpine:latest

WORKDIR /

RUN mkdir /data
RUN addgroup --system otpapp && adduser -S -s /bin/false -G otpapp otpapp

COPY --from=build-env /otpapp /otpapp

RUN chown -R otpapp:otpapp /otpapp
RUN chown -R otpapp:otpapp /data

USER otpapp

EXPOSE 8080

ENTRYPOINT ["/otpapp"]