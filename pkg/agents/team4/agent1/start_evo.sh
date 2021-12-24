#!/bin/bash

# Navigating to root directory
rootdirpath="../../../../"
cd $rootdirpath

# Initialising config files
agentConfigFile="pkg/agents/team4/agent1/agentConfig.json"
bestAgentsFile="pkg/agents/team4/agent1/bestAgents.json"
agentLifeExpectenciesFile="pkg/agents/team4/agent1/agentLifeExpectencies.json"
rm $agentConfigFile $bestAgentsFile
touch $agentConfigFile $bestAgentsFile

degreeOfEquations=1
numberOfBestAgents=4
numberOfAgentsPerSim=20
numberOfIterations=20

python3 pkg/agents/team4/agent1/initaliseConfig.py $agentConfigFile $bestAgentsFile $degreeOfEquations $numberOfBestAgents

for i in $( eval echo {0..$numberOfIterations} )
do
    arr=()
    for j in $( eval echo {0..$numberOfBestAgents} )
    do
        rm logs/*
        make run
        logfile=("logs/*")
        # pass in logfile, num agents, agent_config.json, bestAgent.config, current iteration to python script
        pythonOutput=$(python3 pkg/agents/team4/agent1/checkAgentPerformance.py $logfile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)
        arr+=($pythonOutput)
    done
    printf -v joined '%s,' ${arr[*]}
    echo "[${joined%,}]" > $agentLifeExpectenciesFile
    python3 pkg/agents/team4/agent1/generateNewBestAgents.py $bestAgentsFile $agentLifeExpectenciesFile $degreeOfEquations
done

## Making sure main is unchanged
# cp cmd/backend/main.go cmd/backend/main_copy.go
# sed -i 's/	numOfAgents :=.*/	numOfAgents := []int{10}/g' cmd/backend/main.go

# Generate set of random agents
# for agent in agents (1)
#   create population of only agent
#   (create poulation of all other groups agents + this agent)
#   record average survival rate
# pick out top agent populations
# save list of best agents
# create some hybrid agents
# create some mutated random agents
# return to (1)