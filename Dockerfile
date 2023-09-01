FROM archlinux:latest

ARG uid=1000
ARG gid=1000
ARG USE_FLEXO_CACHE

RUN groupadd -g ${gid} callidus \
 && useradd -u ${uid} -g ${gid} -m -d /domus servus
RUN if [ "z$USE_FLEXO_CACHE" != "z" ]; then \
        echo "Server = http://$(ip r | grep default | awk '{print $3}'):7878/\$repo/os/\$arch" > /etc/pacman.d/mirrorlist; \
    fi; \
    pacman -Syyu --noconfirm \
 && pacman -Sy --noconfirm \
    python \
    make \
    git \
    patch \
    pkgconf \
    which \
    go \
    clang \
    cmake \
    openssl \
    libxml2 \
    mingw-w64-toolchain \
    mingw-w64

USER servus
WORKDIR /domus

RUN git clone https://github.com/tpoechtrager/osxcross \
 && curl -L --output osxcross/tarballs/MacOSX13.3.sdk.tar.xz https://github.com/joseluisq/macosx-sdks/releases/download/13.3/MacOSX13.3.sdk.tar.xz
RUN cd osxcross \
 && UNATTENDED=1 ./build.sh
RUN mkdir scredenziato
COPY . scredenziato

ENV CC_darwin_amd64=/domus/osxcross/target/bin/o64-clang
ENV CXX_darwin_amd64=/domus/osxcross/target/bin/o64-clang++
ENV CC_darwin_arm64=/domus/osxcross/target/bin/aarch64-apple-darwin22.4-clang
ENV CXX_darwin_arm64=/domus/osxcross/target/bin/aarch64-apple-darwin22.4-clang++
ENV CC_windows_amd64=/usr/bin/x86_64-w64-mingw32-gcc
ENV CXX_windows_amd64=/usr/bin/x86_64-w64-mingw32-g++
ENV PATH=$PATH:/domus/osxcross/target/bin

WORKDIR /domus/scredenziato
ENTRYPOINT ["make"]
CMD ["build"]
