FROM postgres:16

RUN apt-get update

ENV build_deps ca-certificates \
  git \
  build-essential \
  libpq-dev \
  postgresql-server-dev-16 \
  curl \
  libreadline6-dev \
  zlib1g-dev

RUN apt-get install -y --no-install-recommends $build_deps pkg-config cmake

WORKDIR /home/postgres

ENV HOME=/home/postgres
ENV PATH=/home/postgres/.cargo/bin:$PATH

RUN chown postgres:postgres /home/postgres

USER postgres

RUN \
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --no-modify-path --profile minimal --default-toolchain 1.70.0 && \
  rustup --version && \
  rustc --version && \
  cargo --version

# pgrx
RUN cargo install cargo-pgrx --version 0.10.2 --locked

RUN cargo pgrx init --pg16 $(which pg_config)

USER root

COPY ./pgx_ulid/Cargo.toml /home/postgres/Cargo.toml
COPY ./pgx_ulid /home/postgres
COPY ./pgx_ulid/postgresql.conf /etc/postgresql/postgresql.conf

RUN cargo pgrx install

RUN chown -R postgres:postgres /home/postgres
RUN chown -R postgres:postgres /usr/share/postgresql/16/extension
RUN chown -R postgres:postgres /usr/lib/postgresql/16/lib

USER postgres

ENV POSTGRES_HOST_AUTH_METHOD=trust
ENV USER=postgres
