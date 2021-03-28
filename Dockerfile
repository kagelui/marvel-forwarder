ARG RELEASE_IMAGE_NAME=alpine
ARG RELEASE_IMAGE_TAG=3.11.5

FROM ${RELEASE_IMAGE_NAME}:${RELEASE_IMAGE_TAG}

LABEL app="marvel-forwarder-api"
LABEL description="marvel API cache"

RUN apk --no-cache add tzdata ca-certificates

COPY ./serverd /

CMD ./serverd
