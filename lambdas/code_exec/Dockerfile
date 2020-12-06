
# Supported code exec languages
FROM public.ecr.aws/bitnami/node:14.15.1-debian-10-r8

FROM public.ecr.aws/lambda/go:1.2020.12.03.17

WORKDIR /build

# Copy code to current directory
COPY go.mod go.sum main.go ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main main.go

ENTRYPOINT [ "/build/main" ]
