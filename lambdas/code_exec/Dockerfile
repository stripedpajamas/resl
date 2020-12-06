
# Supported code exec languages
FROM public.ecr.aws/lambda/nodejs:12

FROM public.ecr.aws/bitnami/golang:latest

WORKDIR /tmp

# Copy code to current directory
COPY go.mod go.sum *.go languages.json ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main *.go

ENTRYPOINT [ "/tmp/main" ]