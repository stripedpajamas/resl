FROM public.ecr.aws/bitnami/golang:latest

# Update system
RUN apt-get update
RUN apt-get -y upgrade

# Install languages
RUN apt-get install -y nodejs
RUN apt-get install -y python3.7

WORKDIR /build

# Copy code to current directory
COPY go.mod go.sum *.go languages.json ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o handler *.go

ENTRYPOINT [ "/build/handler" ]