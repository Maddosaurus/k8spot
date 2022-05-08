# syntax=docker/dockerfile:1

#
# Build container
#
FROM golang:1.18-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
RUN mkdir log

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -o k8spot main.go


#
# Prod container
#
FROM gcr.io/distroless/static AS final

LABEL maintainer="Maddosaurus"
USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /build/k8spot /k8spot
COPY --from=builder --chown=nonroot:nonroot /build/third_party /third_party
COPY --from=builder --chown=nonroot:nonroot /build/log /log


EXPOSE 8080
ENTRYPOINT [ "/k8spot" ]
CMD [ "-flavor minikube" ]