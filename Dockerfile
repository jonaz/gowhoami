FROM alpine:3.13
RUN adduser -u 1001 -D gowhoami
COPY gowhoami /
USER 1001
ENTRYPOINT ["/gowhoami"]
