FROM golang:alpine
LABEL authors="coding_seal"
RUN apk add --no-cache make

WORKDIR /app
COPY . /app
RUN make

ENTRYPOINT ["./task"]