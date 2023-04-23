FROM ubuntu:focal

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get upgrade --yes && DEBIAN_FRONTEND=noninteractive apt-get install --yes \
    tzdata \
    build-essential \
    psmisc \
    net-tools \
    bc \
    tshark \
    iputils-ping \
    iptables \
    vim \
    iperf3 \
    iproute2 \
    libconfig++9v5 \
    libsctp1 \
    pkg-config \
    dh-autoreconf \
    ncurses-dev \
    build-essential \
    libssl-dev \
    libpcap-dev \
    libncurses5-dev \
    libsctp-dev \
    lksctp-tools \
    libprotobuf-dev \
    protobuf-compiler \
    git \
    gpg \ 
    wget \
    curl \
    ninja-build \
    libboost-all-dev \
    libgtest-dev \
    nlohmann-json3-dev \
    libcurl4-openssl-dev \
    rapidjson-dev \
  && rm -rf /var/lib/apt/lists/*

RUN wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc 2>/dev/null | gpg --dearmor - | tee /usr/share/keyrings/kitware-archive-keyring.gpg >/dev/null && \
    echo 'deb [signed-by=/usr/share/keyrings/kitware-archive-keyring.gpg] https://apt.kitware.com/ubuntu/ focal main' | tee /etc/apt/sources.list.d/kitware.list >/dev/null && \
    apt-get update && \
    rm /usr/share/keyrings/kitware-archive-keyring.gpg && \
    apt-get install kitware-archive-keyring && \
    echo 'deb [signed-by=/usr/share/keyrings/kitware-archive-keyring.gpg] https://apt.kitware.com/ubuntu/ focal-rc main' | tee -a /etc/apt/sources.list.d/kitware.list >/dev/null && \ 
    apt-get update && \
    apt-get install cmake -y

WORKDIR /root
RUN git clone https://github.com/nanomsg/nng.git
WORKDIR /root/nng
RUN mkdir build
WORKDIR /root/nng/build
RUN cmake -G Ninja .. && ninja && ninja install && ldconfig

WORKDIR /root
RUN git clone https://github.com/pistacheio/pistache.git
WORKDIR /root/pistache
RUN  git checkout e18ed9baeb2145af6f9ea41246cf48054ffd9907
WORKDIR /root
WORKDIR /root/pistache
RUN mkdir build
WORKDIR /root/pistache/build
RUN cmake .. && make -j8 && make install && ldconfig
WORKDIR /root

# Add sgxTestCApiClient build
COPY ./src /root/sgxtestapiclient
WORKDIR /root/sgxtestapiclient
RUN mkdir build
WORKDIR /root/sgxtestapiclient/build
RUN cmake .. && make -j2

# expose ports
EXPOSE 9090

CMD ["sleep", "infinity"]

