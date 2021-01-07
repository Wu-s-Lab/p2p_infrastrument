FROM ubuntu:18.04

LABEL mantianer="timfan" 




# 记得讲二进制复制到docker同级目录下
COPY IPCertClient /IPCertClient

RUN set -xe \
    && apt update \
    && N2N_VERSION=2.9.0-810 \
    && apt update \
    && apt install --no-install-recommends  -y gnupg iptables lsb-release net-tools wget iproute2 software-properties-common netcat openssh-server docker.io libcap2\
    && add-apt-repository universe \
    && wget https://packages.ntop.org/apt/18.04/all/apt-ntop.deb \
    && apt install /root/apt-ntop.deb \
    && apt update \
    && apt install -y n2n=${N2N_VERSION} \
    && rm -rf /root/apt-ntop.deb /var/lib/apt/lists/* \
    && apt purge -y wget software-properties-common \
    && mkdir /n2n \
    && mkdir /root/.ssh ; mkdir /var/run/sshd \
    && echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config



COPY *.conf /n2n/
COPY *.sh /
COPY authorized_keys /root/.ssh/


EXPOSE 22/udp 22/tcp

RUN chmod +x /start_edge_in_docker.sh ; chmod -R 700 /root/.ssh ; chmod +x /IPCertClient

ENTRYPOINT ["/start_edge_in_docker.sh"]

# ENTRYPOINT [ "/bin/bash" ]



