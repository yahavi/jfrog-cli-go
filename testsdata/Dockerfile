FROM ubuntu:20.04

SHELL ["/bin/bash", "-c"] 

# Environment variables
ENV HOME /root
ENV JAVA_HOME=/usr/lib/jvm/default-java

# Install prerequisites
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -yq zip curl git uuid jq gettext default-jdk golang-1.14-go python3-pip python3-venv --no-install-recommends --no-cache

# Configure Python
RUN ln -s /usr/bin/pip3 /usr/bin/pip
RUN ln -s /usr/bin/python3 /usr/bin/python

# Configure Go
RUN ln -s /usr/lib/go-1.14/bin/go /usr/bin/go

# Install .NET, NuGet, Node.js and npm
RUN curl -sL https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -o packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN curl -sL https://deb.nodesource.com/setup_14.x -o nodesource_setup.sh
RUN apt update
RUN apt install -yq apt-transport-https dotnet-sdk-3.1 nuget nodejs npm --no-install-recommends --no-cache

# Install Maven and Gradle
RUN curl -s "https://get.sdkman.io" | bash
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install maven 3.3.9
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install gradle

# Clean up
RUN apt autoremove
RUN apt clean