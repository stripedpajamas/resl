
FROM public.ecr.aws/lambda/nodejs:12


FROM golang:1.15-alpine

RUN pwd

# Copy run handler folder to current directory
COPY handlers/run_handler /build

WORKDIR /build

# Install packages
RUN go get -v all

# Build run handler
RUN GOOS=linux go build main.go -o main

COPY main /dist

WORKDIR /dist

# Run handler
CMD [ "main" ]
