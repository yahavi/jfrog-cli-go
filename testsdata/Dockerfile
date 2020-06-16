FROM ubuntu:20.04

SHELL ["/bin/bash", "-c"] 

# Environment variables
ENV HOME /root
ENV JAVA_HOME=/usr/lib/jvm/default-java

# Install prerequisites
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -y build-essential zip default-jdk curl python3-pip wget

# Configure python3 to be the default python
RUN ln -s /usr/bin/pip3 /usr/bin/pip
RUN ln -s /usr/bin/python3 /usr/bin/python

# Install .net and Nuget
RUN wget https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN apt update
RUN apt install -y apt-transport-https dotnet-sdk-3.1 nuget

# Install Node JS
RUN curl -sL https://deb.nodesource.com/setup_14.x -o nodesource_setup.sh
RUN apt install -y nodejs npm

# Install Maven and Gradle
RUN curl -s "https://get.sdkman.io" | bash
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install maven 3.3.9
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install gradle
