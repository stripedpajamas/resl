FROM node:14-buster as build-image

# Install aws-lambda-cpp build dependencies
RUN apt-get update && \
    apt-get install -y \
    g++ \
    make \
    cmake \
    unzip \
    libcurl4-openssl-dev

# Copy function code
COPY package*.json /function/
COPY index.js /function/index.js

WORKDIR /function

RUN npm install

# Grab a fresh slim copy of the image to reduce the final size
FROM node:14-buster-slim

# install language runtimes
RUN apt-get update && \
    apt-get install -y \
    python \
    python3

# Include global arg in this stage of the build
ARG FUNCTION_DIR

# Set working directory to function root directory
WORKDIR /function

# Copy in the built dependencies
COPY --from=build-image /function /function

ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /

ENTRYPOINT [ "/entry.sh" ]
CMD ["index.handler"]
