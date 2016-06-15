FROM scratch
#COPY ui/build/ /build
#COPY ui/index.html /
COPY gowhoami /
CMD ["/gowhoami"]
