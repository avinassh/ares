FROM golang:1.9-stretch
ADD build/ /
ENTRYPOINT ["/ares"]