
FROM public.ecr.aws/lambda/nodejs:12


FROM public.ecr.aws/lambda/go:1.2020.12.03.17

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
