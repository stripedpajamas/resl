
# Supported code exec languages
FROM public.ecr.aws/bitnami/node:14.15.1-debian-10-r8

FROM public.ecr.aws/bitnami/golang:latest

WORKDIR /build

# Copy code to current directory
COPY go.mod go.sum *.go languages.json ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main *.go

ENTRYPOINT [ "/build/main" ]