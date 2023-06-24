FROM pactfoundation/pact-cli:latest-multi

RUN apk update  && apk add make git ncurses musl-dev

RUN wget https://golang.org/dl/go1.18.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin