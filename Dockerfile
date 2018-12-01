FROM node:10.12.0-alpine as webappbuilder
COPY ./webapp/ /webapp/
WORKDIR /webapp
RUN yarn install
RUN yarn build

FROM golang:1.11.0 AS gobuilder
COPY ./ /go/src/github.com/andrewbackes/webshell
WORKDIR /go/src/github.com/andrewbackes/webshell
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/webshell cmd/webshell.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=gobuilder /go/src/github.com/andrewbackes/webshell/cmd/webshell /root/webshell
COPY --from=webappbuilder /webapp/build/ /root/webapp/build/
CMD ["/root/webshell"]