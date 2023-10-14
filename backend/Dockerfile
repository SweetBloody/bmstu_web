FROM golang:1.20
WORKDIR /app
COPY ./backend ./
RUN go mod download
RUN #go get github.com/githubnemo/CompileDaemon
#ENTRYPOINT CompileDaemon --build="go build main.go" --command=main
#ENTRYPOINT go run /home/alexey/go/src/github.com/SweetBloody/bmstu_web/backend/main.go
#EXPOSE 5259
CMD ["go", "run", "main.go"]
#RUN CGO_ENABLED=0 GOOS=linux go build -o /
#

