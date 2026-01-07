FROM scratch
ADD omada-to-ntfy /
EXPOSE 8080
CMD ["/omada-to-ntfy"]