# FROM golang:1.16-alpine

# # Set destination for COPY
# WORKDIR /app

# # Download Go modules
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# # Copy the source code. Note the slash at the end, as explained in
# # https://docs.docker.com/engine/reference/builder/#copy
# COPY *.go ./

# # Build
# RUN go build -o /docker-gs-ping

# # This is for documentation purposes only.
# # To actually open the port, runtime parameters
# # must be supplied to the docker command.
# EXPOSE 3000


# # Run
# CMD [ "/docker-gs-ping" ]

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 3000

RUN go build -o main .

CMD ./main