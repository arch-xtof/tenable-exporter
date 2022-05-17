FROM gcr.io/distroless/static

WORKDIR /app

COPY ./tenable-exporter /app
COPY ./mock/* /app/mock/

EXPOSE 9095

CMD ["/app/tenable-exporter"]
