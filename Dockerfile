FROM golang:latest

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

ARG OPENAI_KEY
ARG TG_KEY

ENV OPENAI_KEY=$OPENAI_KEY
ENV TG_KEY=$TG_KEY

RUN echo "OPENAI_KEY=$OPENAI_KEY" >> .env \
    && echo "TG_KEY=$TG_KEY" >> .env

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]
