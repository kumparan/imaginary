FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:builder-vips8.8.4 as builder

ARG IMAGINARY_VERSION="dev"

WORKDIR ${GOPATH}/src/github.com/kumparan/imaginary

# Copy imaginary sources
COPY . .

# Making sure all dependencies are up-to-date
RUN go mod download

# Compile imaginary
RUN CGO_CFLAGS_ALLOW=-Xpreprocessor go test && go build -o ${GOPATH}/bin/imaginary

FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:base-bullseye

ARG SPINNAKER_ID="dev"

LABEL maintainer="aryo.kusumo@kumparan.com" \
      org.label-schema.description="kumparan imaginary" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/kumparan/imaginary" \
      org.label-schema.vcs-url="https://github.com/kumparan/imaginary"

COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /go/src/github.com/kumparan/imaginary/config.yml.* /usr/local/bin/

RUN apt-get update && apt-get install wget -y
RUN wget https://yw-assets.s3-ap-southeast-1.amazonaws.com/Heebo/Heebo-VariableFont_wght.ttf -P /usr/share/fonts/googlefonts