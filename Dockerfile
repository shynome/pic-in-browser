FROM chromedp/headless-shell:128.0.6535.2
RUN apt-get update && \
  apt-get install fonts-noto-cjk
WORKDIR /opt/pic-in-browser
COPY pic-in-browser /opt/pic-in-browser/pic-in-browser
ENTRYPOINT [ "/opt/pic-in-browser/pic-in-browser" ]
