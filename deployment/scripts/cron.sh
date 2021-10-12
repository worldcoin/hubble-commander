#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "General purpose script for maintaining a cron job."
    echo "If used with an exact same command it will be replaced with a new cron schedule."
    echo ""
    echo "NOTE:"
    echo "To remove all cron jobs use:"
    echo "crontab -r"
    echo ""
    echo "WARNING:"
    echo "  - The cron job will be added to the crontab file of the user who calls it!"
    echo "  - If the script returns a warning \"crontab: no crontab for USER\", please ignore it!"
    echo ""
    echo "USAGE:"
    echo "Script requires 2 arguments:"
    echo "  1. cron schedule expression"
    echo "  2. command"
    echo ""
    echo "Example usage:"
    echo "bash $0 \"5 */1 * * *\" \"bash /home/user/script.sh argument1 argument 2\""
    exit 0
fi

# Delete a crontab entry with the same command
crontab  -l | grep -v "$2"  | crontab -

# Append a new command with a cron schedule
(crontab -l ; echo "$1" "$2") | crontab -
