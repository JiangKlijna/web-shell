FROM golang:1.21 AS build-env
WORKDIR /app
COPY  . /app
RUN useradd -u 1000 webshell
RUN go mod tidy
RUN go mod vendor
RUN make gen
RUN make

FROM ubuntu:22.04
# Upgrade system & install packages
ENV DEBIAN_FRONTEND=noninteractive
ENV TERM=linux
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y -o Dpkg::Options::=--force-confdef -o Dpkg::Options::=--force-confnew vim-tiny bash curl wget git telnet

# Add kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x kubectl && mv kubectl /usr/local/bin/

# Add web-shell
COPY --from=build-env /app/web-shell /web-shell
COPY --from=build-env /etc/passwd /etc/passwd
USER webshell

# Set env vars
ENV HOST=0.0.0.0
ENV PORT=2020
ENV USER=webshell
ENV PASSWORD=webshell

ENTRYPOINT ["/web-shell -s -H $HOST -P $PORT -u $USER -p $PASSWORD"]
