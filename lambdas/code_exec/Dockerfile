
# Supported code exec languages
FROM public.ecr.aws/lambda/nodejs:12

FROM public.ecr.aws/bitnami/golang:latest

COPY go.mod go.sum *.go languages.json ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o handler *.go

ENTRYPOINT [ "handler" ]