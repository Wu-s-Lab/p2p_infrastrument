#!/bin/bash

EDGE_CONF_DIR=/edge
EDGE_CONF_FILE=edge.conf


echo "The host system is $HOST"

# Start the sshd process
/usr/sbin/sshd -f /n2n/sshd.conf
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start sshd: $status"
  exit $status
fi


# Start the edge conf reqest process

if [ ! -f "$EDGE_CONF_DIR/$EDGE_CONF_FILE" ];then

        
    echo "Getting the conf req of $NAME , the phone is $PHONE "

    if [ ! -f "$EDGE_CONF_DIR/client.crt" ];then
        /IPCertClient -c init -s $IPCERT_URL  -f $EDGE_CONF_DIR -n $NAME -p $PHONE
        sleep 3
    fi
    /IPCertClient -c getEdgeConf -f  $EDGE_CONF_DIR/client.crt -e $EDGE_CONF_DIR/$EDGE_CONF_FILE -s $IPCERT_URL
    
    
    status=$?
    if [ $status -ne 0 ]; then
    echo "Failed to request edge config: $status"
    exit $status
    fi

fi

# Start the edge process if the HOST is linux
if [ $HOST = 'Linux' ]; then
    echo "Starting edge node"
    edge $EDGE_CONF_DIR/$EDGE_CONF_FILE
    status=$?
    if [ $status -ne 0 ]; then
    echo "Failed to start n2n process: $status"
    exit $status
    fi
fi





# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds

while sleep 60; do
  ps aux |grep sshd |grep -q -v grep
  PROCESS_1_STATUS=$?

  # If the greps above find anything, they exit with 0 status
  # If they are not both 0, then something is wrong
  if [ $PROCESS_1_STATUS -ne 0 ]; then
    echo "sshd  processes has already exited."
    exit 1
  fi
  
  if [ $HOST = 'Linux' ]; then
    ps aux |grep edge |grep -q -v grep
    PROCESS_2_STATUS=$?

    if [ $PROCESS_2_STATUS -ne 0 ]; then
        echo "edge processes has already exited."
        exit 1
    fi
  fi

done
