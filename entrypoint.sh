#!/bin/bash

cd /home/aislab/agent/abusim-goabu-agent/
source /opt/ros/humble/setup.bash && source /home/aislab/aburos/aburos_msgs/install/setup.bash && source /home/aislab/rosetta/install/setup.bash

./abusim-goabu-agent $1