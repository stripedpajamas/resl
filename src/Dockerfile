
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/bitnami/golang:latest

# Copy run handler folder to current directory
COPY handlers/ /build

WORKDIR /build

# Install packages
RUN go get -v all

# Build run handler
RUN GOOS=linux go build -o main run_handler/main.go

# Run handler
CMD [ "main" ]
