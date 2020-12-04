
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/bitnami/golang:latest

WORKDIR /build

# Copy run handler folder to current directory
COPY go.mod .
COPY go.sum .

# Install packages
RUN go get -v all

# Build run handler
RUN GOOS=linux go build -o main main.go

# Run handler
ENTRYPOINT [ "/build/main" ]
