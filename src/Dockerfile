
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/bitnami/golang:1.14.12-debian-10-r11

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
