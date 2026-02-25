#!/bin/bash

cd /home/aislab/agent/abusim-goabu-agent/
source /opt/ros/humble/setup.bash && source /home/aislab/aburos/aburos_msgs/install/setup.bash && source /home/aislab/goROSetta/goROSetta_msgs/install/setup.bash

./abusim-goabu-agent $1