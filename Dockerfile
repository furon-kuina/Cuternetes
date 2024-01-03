FROM docker:dind

WORKDIR /app/src/cmd/cutelet

RUN wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin"
ENV GOBIN="/usr/local/go/bin"
RUN echo $PATH
RUN echo $GOBIN
RUN go install github.com/cosmtrek/air@latest
