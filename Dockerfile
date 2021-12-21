FROM golang:1.17-alpine AS build

ARG FIREBASE_ADMINSDK

# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/application
COPY . .
ENV GOPATH="/go/src"
RUN GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/application/main .

RUN echo $FIREBASE_ADMINSDK > aumsu-portal-firebase-adminsdk-5sajn-e6d3adfd5a.json
RUN mkdir -p /var/www/images/messages
RUN mkdir -p /var/www/images/avatars
EXPOSE 8010
ENTRYPOINT  ["./main"]
