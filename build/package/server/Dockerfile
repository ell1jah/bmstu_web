FROM golang:1.20.4-alpine AS build_stage
COPY ./images /go/src/app/images
COPY ./cmd /go/src/app/cmd
COPY ./model /go/src/app/model
COPY ./internal /go/src/app/internal
COPY ./go.* /go/src/app/
WORKDIR /go/src/app
RUN  go mod download
RUN  go install ./cmd/main.go

FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/bin/main /app_binary/
COPY --from=build_stage /go/src/app/images /app_binary/images/
RUN chmod +x main
ENTRYPOINT ./main

EXPOSE 8080/tcp
CMD [ "main" ]
