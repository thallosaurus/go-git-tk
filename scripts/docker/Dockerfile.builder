FROM fedora:40
RUN dnf install -y golang @development-tools git openssh openssh-server
WORKDIR /app
COPY . .
RUN make all

RUN export GOPATH=$(go env GOPATH)
RUN export PATH=$PATH:$GOPATH/bin
RUN mkdir -p /var/local/git-tk/repos
RUN useradd -M -d /var/local/git-tk git
RUN chown git:git /var/local/git-tk/repos

#RUN cp /app/scripts/sshd.d/92-git-config.conf /etc/ssh/sshd_config.d/02-git-config.conf

RUN ssh-keygen -i /root/.ssh/id_ed25519

WORKDIR /var/local/git-tk
CMD ["/usr/sbin/sshd", "-D"]