# SPDX-License-Identifier: ice License 1.0

FROM golang:latest AS build
ARG SERVICE_NAME
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/
COPY . /app/

ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN env SERVICE_NAME=$SERVICE_NAME make dockerfile
RUN cp cmd/$SERVICE_NAME/bin bin

FROM gcr.io/distroless/base-debian11:latest
ARG TARGETOS
ARG TARGETARCH
ARG PORT=443
LABEL os=$TARGETOS
LABEL arch=$TARGETARCH
COPY --from=build /app/bin app
#You might need to expose more ports. Just add more separated by space
#I.E. EXPOSE 8080 8081 8082 8083
EXPOSE $PORT
ENTRYPOINT ["/app"]
