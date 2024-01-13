# syntax=docker/dockerfile:1

FROM golang:1.21.0

ENV GO111MODULE=on
# Set destination for COPY
WORKDIR /build
ADD . .
RUN go env -w GO111MODULE=on
# Download Go modules
# COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
# COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix tgo -o app .



# COPY --chown=admin:admin --from=build /build/app /usr/bin/app
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080
ENTRYPOINT [ "/build/app"]
# Run
# CMD ["/trackprosto"]