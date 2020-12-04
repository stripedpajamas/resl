
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/bitnami/golang:latest

# Copy run handler folder to current directory
COPY handlers/ /build

WORKDIR /build

# Install packages
RUN go get -v all

# Build run handler
RUN GOOS=linux go build -o ./dist/main run_handler/main.go

RUN cd ./dist

RUN pwd

RUN ls -alt

# Run handler
CMD [ "main" ]
