FROM alpine:latest
MAINTAINER Xin Wang <hsimwong1202@gmail.com>

COPY ./build/tx2-k8s-device-plugin /root/tx2-k8s-device-plugin
RUN touch /home/gpu-device

RUN mkdir /var/lib/kubelet/
RUN mkdir /var/lib/kubelet/device-plugins/

RUN touch /var/lib/kubelet/device-plugins/csu-gpu.sock
CMD ["/root/tx2-k8s-device-plugin"]
