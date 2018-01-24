FROM golang:alpine
WORKDIR /app
ADD ./app /app/app
RUN cd /app
EXPOSE 12001
# && go build -o goapp
ENTRYPOINT ["./app"]