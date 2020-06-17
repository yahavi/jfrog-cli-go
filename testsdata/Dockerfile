FROM ubuntu:20.04

SHELL ["/bin/bash", "-c"] 

# Environment variables
ENV HOME /root
ENV JAVA_HOME /root/.sdkman/candidates/java/current
ENV PATH /root/.sdkman/candidates/java/current/bin:/root/.sdkman/candidates/maven/current/bin:/root/.sdkman/candidates/gradle/current/bin:/usr/lib/go-1.14/bin:${PATH}
ENV M2_HOME /root/.sdkman/candidates/maven/current

# Install prerequisites
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -yq zip unzip curl git uuid jq gettext golang-1.14-go python3-pip python3-venv

# Configure Python
RUN ln -s /usr/bin/pip3 /usr/bin/pip
RUN ln -s /usr/bin/python3 /usr/bin/python

# Install .NET, NuGet, Node.js and npm
RUN curl -sL https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -o packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN curl -sL https://deb.nodesource.com/setup_12.x -o nodesource_setup.sh
RUN bash nodesource_setup.sh
RUN apt update
RUN apt install -yq apt-transport-https dotnet-sdk-3.1 nuget nodejs npm

# Install Java, Maven and Gradle
RUN curl -s "https://get.sdkman.io" | bash
RUN source "/root/.sdkman/bin/sdkman-init.sh" && sdk install java `sdk list java | grep -E "11.*-zulu" | head -1 | awk '{print $NF}'` && java -version \
    && sdk install maven 3.3.9 \
    && sdk install gradle \
    && sdk flush archives

# Clean up
RUN apt autoremove
RUN apt clean
