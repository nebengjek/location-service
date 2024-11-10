FROM golang:1.23.3-alpine

LABEL maintainer="farid.wicak@gmail.com"

WORKDIR /usr/src/app

# Update package
RUN apk add --update --no-cache --virtual .build-dev build-base git

COPY . .

RUN make install \
  && make build

# Expose port
EXPOSE 9000

# Run application
CMD ["make", "start"]

