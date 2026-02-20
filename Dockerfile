# Create the building image for compiling
FROM lucagemolotto/aburos-container AS build

WORKDIR /home/aislab
RUN mkdir ./agent

COPY ./abusim-core /home/aislab/abusim-core/
WORKDIR /home/aislab/abusim-core/schema
RUN go mod download -x

WORKDIR /home/aislab/agent/abusim-goabu-agent
COPY --chown=aislab:aislab ./abusim-goabu-agent ./abusim-goabu-agent/entrypoint.sh ./
RUN chmod +x entrypoint.sh
RUN go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/aburos \
 && go mod edit -dropreplace=github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/aburos=../../aburos \
 && go mod edit -replace=github.com/Autonomous-Systems-Laboratory-UNIUD/abusim-core/schema=../../abusim-core/schema
RUN go mod tidy
RUN go mod download -x
RUN cat go.mod

RUN source /opt/ros/humble/setup.bash && source /home/aislab/aburos/aburos_msgs/install/setup.bash && source /home/aislab/rosetta/install/setup.bash && go build

ENTRYPOINT [ "/home/aislab/agent/abusim-goabu-agent/entrypoint.sh" ]
