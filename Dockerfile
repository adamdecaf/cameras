FROM ubuntu:15.10

# install opencv things
RUN apt-get update # && \
    # apt-get install

# copy conf files over
ADD haarcascade_frontalface_alt.xml /opt/cameras/haarcascade_frontalface_alt.xml
ADD cameras /opt/cameras/cameras
