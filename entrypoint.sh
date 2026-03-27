#!/bin/bash

cd /home/aislab/agent/abusim-goabu-agent/
source /opt/ros/humble/setup.bash && source /home/aislab/aburos/aburos_msgs/install/setup.bash && source /home/aislab/goROSetta/goROSetta_msgs/install/setup.bash
echo "export ROS_LOCALHOST_ONLY=1" >> ~/.bashrc

sed \
  -e "s/AGENT_ID_PLACEHOLDER/${HOSTNAME}/" \
  /home/aislab/agent/abusim-goabu-agent/zenoh_bridge.json5.template > /home/aislab/agent/abusim-goabu-agent/bridge.json5

export ROS_LOCALHOST_ONLY=1 && zenoh-bridge-ros2dds -c /home/aislab/agent/abusim-goabu-agent/bridge.json5 &

export ROS_LOCALHOST_ONLY=1 && ./abusim-goabu-agent $1