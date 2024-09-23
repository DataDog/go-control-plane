FROM golang:1.22@sha256:4594271250150c1a322ed749abfd218e1a8c6eb1ade90872e325a664412e2037

WORKDIR /go/src/dyplomat
COPY . /go/src/dyplomat

RUN go install -v ./...
CMD ["dyplomat"]
