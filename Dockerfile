FROM ubuntu:latest
MAINTAINER microbusinesses.inc@gmail.com
ADD AddressService /
EXPOSE 80
CMD ["/AddressService"]
