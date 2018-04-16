#!/bin/bash


# Run this script to configure the environment where the logserver
# will be executed.
#   - notes:
#       This script is not versatile.  A production quality
#       version would need to check OS/architecture/tools
#       in order to properly chroot the logserver and
#       properly enforce permissions.

NO_CHROOT=1

NEW_GROUP="apptiologserver"
SCP_DIR="./scpdir"
MAIN_APP_LOG_DIR="/home/bbright/mainapplog/"
MAIN_APP_LOG_FILE="mainapp.log"
MAIN_LOG_FULL=$MAIN_APP_LOG_DIR$MAIN_APP_LOG_FILE
LOG_USER="logsserveruser"

# Changes group settings of the main app log and creates a new user
# that will run the logserver
groupadd $NEW_GROUP
mv $SCP_DIR/* $MAIN_APP_LOG_DIR
useradd $LOG_USER -g $NEW_GROUP
chown -R :$NEW_GROUP $MAIN_APP_LOG_DIR
chmod -R g+r $MAIN_APP_LOG_DIR

# Add cron rule to change group of files in the main app log directory,
#   in case new log files are created by the main app
crontab -l > tempcron
echo "00 00 * * * chown -R :$NEW_GROUP $MAIN_APP_LOG_DIR && chmod -R g+r $MAIN_APP_LOG_DIR " >> tempcron
crontab tempcron
rm tempcron


if [ $NO_CHROOT -eq 1 ]
then
    echo "chrooting not enabled.  exiting"
    exit 0
fi

# Now create a chroot jail in order to limit what an attacker can
# do after gaining shell access.  
#    - this is 

cd  $MAIN_APP_LOG_DIR
mkdir bin dev etc etc/pam.d home lib lib/security var var/log usr usr/bin

# Copy only the libraries needed to run chroot
# Below is my first attempt 
ldd ./logserver | awk '{print $1}' | xargs -I{} cp "/lib/{}" /home/bbright/mainapplog/lib/


# Finally, run the chroot
chroot . /logserver
