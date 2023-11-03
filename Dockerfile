FROM node:21.1-bullseye-slim as angular_build

COPY frontend /frontend
WORKDIR frontend
# First install all dependencies, because we need dev dependencies for building
RUN npm ci

ENV NODE_ENV production
# Make two layers, because dependencies don't change as often as application
RUN npm run build --progress=false

FROM golang:1.19-alpine as build

COPY backend /backend
WORKDIR /backend

RUN go mod tidy && \
    go mod download

# Disable crosscompiling
ENV CGO_ENABLED=0
# Compile Linux only
ENV GOOS=linux
# Build the binary with debug information removed
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o /main .

# Start with a scratch (no layers)
FROM scratch

# Copy the ca certificates from the build image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=angular_build /frontend/dist/frontend /frontend/dist/frontend

WORKDIR backend
# Copy our static linked library
COPY --from=build main main

EXPOSE 8080

ENTRYPOINT ["./main"]
