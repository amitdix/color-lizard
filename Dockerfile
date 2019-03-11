FROM golang:latest

WORKDIR /go/src/color-lizard

COPY . .
ARG drone_tag=1.0.0
ARG drone_commit
ENV VERSION=${drone_tag}
ENV IMAGE=${drone_commit}

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o colorlizard ./cmd/colorlizard.go

ENV GIN_MODE=release \
  PORT=8080
RUN ls -ltr
RUN pwd

EXPOSE 8080

CMD ./colorlizard
