
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/bitnami/golang:latest

WORKDIR /build

# Copy run handler folder to current directory
COPY go.mod go.sum main.go ./

# Install dependencies
RUN go mod download

# Build code
RUN GOOS=linux go build -o main main.go

ENTRYPOINT [ "/build/main" ]
