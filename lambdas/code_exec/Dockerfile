FROM public.ecr.aws/bitnami/golang:latest

RUN apt-get update
RUN apt-get install -y nodejs

WORKDIR /build

# Copy code to current directory
COPY go.mod go.sum *.go languages.json ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main *.go

ENTRYPOINT [ "/build/main" ]