FROM ubuntu:latest
RUN useradd -u 1001 gowhoami

FROM scratch
COPY --from=0 /etc/passwd /etc/passwd
COPY gowhoami /
USER gowhoami
ENTRYPOINT ["/gowhoami"]
