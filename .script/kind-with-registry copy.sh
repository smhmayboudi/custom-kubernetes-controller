#!/bin/sh
set -o errexit

export CHANNELS=development
export DEFAULT_CHANNEL=development
export USERNAME=smhmayboudi
export VERSION=0.0.1
export IMAGE_TAG_BASE=localhost:5000/$USERNAME/custom-kubernetes-controller
export IMG=$IMAGE_TAG_BASE:v$VERSION

