FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:builder as builder

ARG IMAGINARY_VERSION="dev"

WORKDIR ${GOPATH}/src/github.com/kumparan/imaginary

# Copy imaginary sources
COPY . .

# Making sure all dependencies are up-to-date
RUN rm -rf vendor && dep init && dep ensure

# # Compile imaginary
RUN make install && make build

FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:base

ARG IMAGINARY_VERSION

LABEL maintainer="tomas@aparicio.me" \
      org.label-schema.description="Fast, simple, scalable HTTP microservice for high-level image processing with first-class Docker support" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/h2non/imaginary" \
      org.label-schema.vcs-url="https://github.com/h2non/imaginary" \
      org.label-schema.version="${IMAGINARY_VERSION}"

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Server port to listen
ENV PORT 9000

# Run the entrypoint command by default when the container starts.
ENTRYPOINT ["/usr/local/bin/imaginary"]

# Expose the server TCP port
EXPOSE ${PORT}