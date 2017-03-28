FROM scratch
#COPY ui/build/ /build
#COPY ui/index.html /
COPY gowhoami /
ENTRYPOINT ["/gowhoami"]
