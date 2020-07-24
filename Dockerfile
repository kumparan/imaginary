FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:builder-ubuntu-vips8.8.4-libjpeg9-libpng16 as builder

ARG IMAGINARY_VERSION="dev"

RUN export GOPATH=/tmp/go

WORKDIR ${GOPATH}/src/github.com/kumparan/imaginary

# Copy imaginary sources
COPY . .

# Cache go modules
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
COPY deps.sh .

RUN go mod download

RUN chmod +x ./deps.sh
RUN ./deps.sh

# Compile imaginary
RUN CGO_CFLAGS_ALLOW=-Xpreprocessor go test && go build -o ${GOPATH}/bin/imaginary

FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:base-ubuntu-libjpeg9-libpng16

ARG SPINNAKER_ID="dev"

LABEL maintainer="aryo.kusumo@kumparan.com" \
      org.label-schema.description="kumparan imaginary" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/kumparan/imaginary" \
      org.label-schema.vcs-url="https://github.com/kumparan/imaginary"

COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /go/src/github.com/kumparan/imaginary/config.yml.* /usr/local/bin/