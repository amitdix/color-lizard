FROM docker.target.com/sdp/goalpine:stable

WORKDIR /go/src/git.target.com/StoreDataMovement/color-lizard

COPY . .
ARG drone_tag=1.0.0
ARG drone_commit
ENV VERSION=${drone_tag}
ENV IMAGE=${drone_commit}

RUN go build -ldflags "-X main.version=${VERSION} -X main.image=${IMAGE}" ./cmd/color-lizard.go

FROM docker.target.com/tap/alpine-certs

ENV GIN_MODE=release \
  PORT=8080

WORKDIR /cmd/
RUN ls

EXPOSE 8080
COPY --from=0 /go/src/git.target.com/StoreDataMovement/color-lizard/color-lizard .
RUN ls

CMD ./color-lizard