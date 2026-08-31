# Create the building image for compiling
FROM lucagemolotto/aburos-container AS build

WORKDIR /home/aislab
RUN mkdir ./agent

## Zenoh configuration
RUN echo "export ROS_LOCALHOST_ONLY=1" >> ~/.bashrc
RUN sudo apt install iproute2 -y
#RUN sudo ip l set lo multicast on
RUN curl -L https://download.eclipse.org/zenoh/debian-repo/zenoh-public-key | sudo gpg --dearmor --yes --output /etc/apt/keyrings/zenoh-public-key.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/zenoh-public-key.gpg] https://download.eclipse.org/zenoh/debian-repo/ /" | sudo tee -a /etc/apt/sources.list > /dev/null \
    && sudo apt update
RUN sudo apt install -y zenoh-bridge-ros2dds=1.8.0

COPY ./abusim-core /home/aislab/abusim-core/
WORKDIR /home/aislab/abusim-core/schema
RUN go mod download -x

WORKDIR /home/aislab/agent/abusim-goabu-agent
COPY --chown=aislab:aislab ./abusim-goabu-agent ./abusim-goabu-agent/entrypoint.sh ./
RUN chmod +x entrypoint.sh
#ENV FASTRTPS_DEFAULT_PROFILES_FILE=/home/aislab/agent/abusim-goabu-agent/fastdds_profile.xml
#ENV RMW_IMPLEMENTATION=rmw_fastrtps_cpp
RUN sudo apt install -y ros-humble-rmw-cyclonedds-cpp
ENV RMW_IMPLEMENTATION=rmw_cyclonedds_cpp
ENV GOTRACEBACK=all
ENV GOFLAGS=-trimpath
RUN go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/aburos \
 && go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema \
 && go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/goROSetta/ROSetta \
 && go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/goROSetta/goMavUtil \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/aburos=../../aburos \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema=../../abusim-core/schema \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/goROSetta/ROSetta=../../goROSetta/ROSetta \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/goROSetta/goMavUtil=../../goROSetta/goMavUtil
RUN go mod tidy
RUN go mod download -x
RUN cat go.mod

RUN source /opt/ros/humble/setup.bash && source /home/aislab/aburos/aburos_msgs/install/setup.bash && source /home/aislab/rosetta/install/setup.bash && source /home/aislab/goROSetta/goROSetta_msgs/install/setup.bash && go build

ENTRYPOINT [ "/home/aislab/agent/abusim-goabu-agent/entrypoint.sh" ]
