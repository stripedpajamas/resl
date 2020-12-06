
# Supported code exec languages
FROM public.ecr.aws/lambda/nodejs:12

FROM public.ecr.aws/bitnami/golang:latest

WORKDIR /build

# Copy code to current directory
COPY go.mod go.sum *.go ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main *.go

ENTRYPOINT [ "/build/main" ]