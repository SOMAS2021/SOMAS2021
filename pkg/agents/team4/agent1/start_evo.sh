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

numberOfHealthLevels=4
numberOfBestAgents=4
numberOfAgentsPerSim=20
numberOfIterations=19
numberOfRuns=3

# Generate set of agents with 0 parameters
python3 pkg/agents/team4/agent1/initaliseConfig.py $agentConfigFile $bestAgentsFile $numberOfHealthLevels $numberOfBestAgents

for i in $( eval echo {0..$numberOfIterations} )
do
    echo "ITERATION " $i
    echo ""
    arr=()
    for j in $( eval echo {0..$numberOfBestAgents} )
    do
        echo "  Getting average performance of agent " $j
        individualResults="0.0"
        for k in $( eval echo {1..$numberOfRuns} )
        do
            rm logs/*
            # create population of only agent
            # (could be changed in future for other groups agents + this agent)
            make run
            logfile=("logs/*")
            # pass in logfile, num agents, agent_config.json, bestAgent.config, current iteration to python script
            pythonOutput=$(python3 pkg/agents/team4/agent1/checkAgentPerformance.py $logfile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)
            # record average survival rate
            # echo "      Run " $k " - " $pythonOutput 
            # individualResults="$individualResults + $pythonOutput" | python3
            individualResults=`echo $individualResults+$pythonOutput | bc`
        done
        averageResult=`echo $individualResults/$numberOfRuns | bc -l` 
        # echo "  Average perfomance " $individualResults 
        arr+=($averageResult)
    done
    printf -v joined '%s,' ${arr[*]}
    echo "[${joined%,}]" > $agentLifeExpectenciesFile
    # generate new set of best agents generated from previous perfomance 
    python3 pkg/agents/team4/agent1/generateNewBestAgents.py $bestAgentsFile $agentLifeExpectenciesFile $numberOfHealthLevels
    echo "------------------------------------------"
done

# Initialise set of agents
# for agents in all agents
#   create population of only agent
#   (could be changed in future for other groups agents + this agent)
#   record average survival rate
# pick out top agent populations
# save list of best agents
# create some hybrid agents
# create some mutated random agents
# return to (1)