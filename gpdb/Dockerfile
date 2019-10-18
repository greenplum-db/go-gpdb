FROM centos:7

COPY gpdb config.yml greenplum-db-6.0.1-rhel7-x86_64.rpm /root/

# Set up SSH
RUN yum install -y openssh-server openssh-clients sudo sshpass && \
    ssh-keygen -N '' -f /etc/ssh/ssh_host_rsa_key && \
    echo '/sbin/sshd' >> /root/.bashrc && \
    echo root:root | chpasswd && \
    yum clean all

# Set up gpadmin
RUN useradd gpadmin && \
    echo -e "root\nroot" | passwd gpadmin && \
    echo 'gpadmin  ALL=(ALL:ALL) ALL' >> /etc/sudoers && \
    chmod 770 /usr/local/ && \
    chgrp gpadmin /usr/local/

# Install gpdb 6.0.1
RUN mkdir -p /home/gpadmin/gpdbinstall/download && \
    mv /root/greenplum-db-6.0.1-rhel7-x86_64.rpm /home/gpadmin/gpdbinstall/download && \
    mv /root/gpdb /root/config.yml /home/gpadmin && \
    chown -R gpadmin:gpadmin /home/gpadmin/*

CMD su - gpadmin
#echo $HOSTNAME > /home/gpadmin/hostfile
#/sbin/sshd
#/home/gpadmin/gpdb install -v 6.0.1