FROM golang:1.13
ENV GO111MODULE=on
WORKDIR /sensorapi/
COPY . .
RUN ls
RUN GOOS=linux GOARCH=amd64 go build ./main.go
RUN mv ./main ./sapi
CMD ./sapi
