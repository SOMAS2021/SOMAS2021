#!/bin/bash

# Navigating to root directory
rootdirpath="../../../../"
cd $rootdirpath

# Initialising config file
configFile="pkg/agents/team4/agent1/config.json"
rm $configFile
touch $configFile
python3 pkg/agents/team4/agent1/initaliseConfig.py $configFile

## Making sure main is unchanged
# cp cmd/backend/main.go cmd/backend/main_copy.go
# sed -i 's/	numOfAgents :=.*/	numOfAgents := []int{10}/g' cmd/backend/main.go

# Running simulation once
for i in {0..3}
do
    echo "iteration: "$i
    rm logs/*
    make run
    logfile=("logs/*")
    python3 pkg/agents/team4/agent1/checkAgentPerformance.py $logfile 10 $configFile
    echo "---------------------------------------------------------------"
done