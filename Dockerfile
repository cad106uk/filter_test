FROM golang:alpine
WORKDIR /app
ADD . /app
RUN cd /app && go build -o filtered_test
ENTRYPOINT ./filtered_test
