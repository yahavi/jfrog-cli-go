FROM ubuntu:20.04

SHELL ["/bin/bash", "-c"] 

# Environment variables
ENV HOME /root
ENV JAVA_HOME=/usr/lib/jvm/default-java

# Install prerequisites
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -y wget zip curl git uuid jq gettext default-jdk python3-pip python3-venv

# Configure python3 and pip3 to be the default Python and Pip
RUN ln -s /usr/bin/pip3 /usr/bin/pip
RUN ln -s /usr/bin/python3 /usr/bin/python

# Install .NET and NuGet
RUN wget https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN apt update
RUN apt install -y apt-transport-https dotnet-sdk-3.1 nuget

# Install Node.js and npm
RUN curl -sL https://deb.nodesource.com/setup_14.x -o nodesource_setup.sh
RUN apt install -y nodejs npm

# Install Maven and Gradle
RUN curl -s "https://get.sdkman.io" | bash
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install maven 3.3.9
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install gradle
