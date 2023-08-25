FROM golang:latest

COPY . /github.com/andy-ahmedov/telegram_bot/

WORKDIR /github.com/andy-ahmedov/telegram_bot/

# RUN apt-get update &&\
# 	apt-get -y install postgresql-client
# # make wait-for-postgres.sh executable
# RUN chmod +x wait-for-postgres.sh
# RUN go mod download
# RUN go build -o ./bin/bot cmd/bot/main.go
# CMD ["./bin/bot"]

#--------------------------
RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /github.com/andy-ahmedov/telegram_bot/.bin/bot .
COPY --from=0 /github.com/andy-ahmedov/telegram_bot/configs configs/
COPY --from=0 /github.com/andy-ahmedov/telegram_bot/wait-for-postgres.sh .
COPY --from=0 /github.com/andy-ahmedov/telegram_bot/fiz_lica.xml .
RUN apk update &&\
	apk add postgresql-client
RUN chmod +x wait-for-postgres.sh
RUN apk add libc6-compat
CMD ["./bot"]
#--------------------------


