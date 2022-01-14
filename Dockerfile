FROM golang:1.16-buster
RUN apt-get update && apt-get install -y --no-install-recommends unzip openjdk-11-jdk
WORKDIR /android
ENV ANDROID_HOME=/android \
    ANDROID_VERSION=29 \
    ANDROID_BUILD_TOOLS_VERSION=30.0.2 \
    NDK_VERSION=23.0.7599858
RUN curl -O https://dl.google.com/android/repository/commandlinetools-linux-7583922_latest.zip && \
    unzip commandlinetools-linux-7583922_latest.zip && \
    rm -f commandlinetools-linux-7583922_latest.zip
RUN yes | ./cmdline-tools/bin/sdkmanager --sdk_root=$ANDROID_HOME --licenses
RUN ./cmdline-tools/bin/sdkmanager --sdk_root=$ANDROID_HOME \
    "build-tools;${ANDROID_BUILD_TOOLS_VERSION}" \
    "platforms;android-${ANDROID_VERSION}" \
    "platform-tools" \
    "ndk;${NDK_VERSION}"
RUN ln -sf $ANDROID_HOME/ndk/${NDK_VERSION} $ANDROID_HOME/ndk-bundle
ENV GO111MODULE=off
RUN go get -v golang.org/x/mobile/cmd/gomobile
RUN gomobile init
WORKDIR /go/src/github.com/caiguanhao/adboverssh
COPY . .
RUN go get -v
RUN cd mobile && gomobile bind -target=android/arm -v -o /adboverssh-arm.aar .
