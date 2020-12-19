#!/bin/bash


LAST_COMMIT=$(git rev-parse HEAD)

pushd lambdas/slack_listener
echo "Updating slack_listener..."
go get "github.com/stripedpajamas/resl/models@${LAST_COMMIT}"
go get "github.com/stripedpajamas/resl/slack@${LAST_COMMIT}"
echo "Building slack_listener..."
go build
popd

pushd lambdas/slack_responder
echo "Updating slack_responder..."
go get "github.com/stripedpajamas/resl/models@${LAST_COMMIT}"
go get "github.com/stripedpajamas/resl/slack@${LAST_COMMIT}"
echo "Building slack_responder..."
go build
popd

echo "Updated lambda go.mods with latest models and slack module commits"

