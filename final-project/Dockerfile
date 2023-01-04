FROM docker.io/bitnami/golang:1.19 AS builder
ADD ./ /app

WORKDIR /app
RUN rm go.mod && go mod init github.com/calvarado2004/go-concurrency && go mod tidy && go get github.com/alexedwards/scs/v2  && go get github.com/xhit/go-simple-mail/v2 && go get github.com/jackc/pgconn && go get github.com/jackc/pgx/v4 && go get github.com/go-chi/chi/v5 && go get github.com/bwmarrin/go-alone && go get github.com/phpdave11/gofpdf && go get github.com/phpdave11/gofpdf/contrib/gofpdi

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mailconcurrent cmd/web/*.go

FROM busybox:stable
ENV APP_HOME /app
RUN adduser 1001 -D -h $APP_HOME && mkdir -p $APP_HOME && chown 1001:1001 $APP_HOME
USER 1001
WORKDIR $APP_HOME
COPY ./cmd/web/templates /app/cmd/web/templates/
COPY --chown=0:0 --from=builder /app/mailconcurrent ./
EXPOSE 8080
CMD ["./mailconcurrent"]