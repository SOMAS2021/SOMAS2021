#!/bin/bash

# Navigating to root directory
rootdirpath="../../../../"
cd $rootdirpath

# Initialising config files
agentConfigFile="pkg/agents/team4/agent1/agentConfig.json"
bestAgentsFile="pkg/agents/team4/agent1/bestAgents.json"
agentLifeExpectanciesFile="pkg/agents/team4/agent1/agentLifeExpectancies.json"
agentDeathRateFile="pkg/agents/team4/agent1/agentDeathRate.json"
rm $agentConfigFile $bestAgentsFile
touch $agentConfigFile $bestAgentsFile

numberOfHealthLevels=4
numberOfBestAgents=5
numberOfAgentsPerSim=16
numberOfIterations=1
numberOfRuns=1

# Generate set of agents with 0 parameters
python3 pkg/agents/team4/agent1/initaliseConfig.py $agentConfigFile $bestAgentsFile $numberOfHealthLevels $numberOfBestAgents

for i in $( eval echo {1..$numberOfIterations} )
do
    echo "ITERATION " $i
    echo ""
    arrLifeExp=()
    arrDeathRate=()
    for j in $( eval echo {1..$numberOfBestAgents} )
    do
        echo "  Getting average performance of agent " $j
        averageLifeExpectancy="0.0"
        averageDeathRate="0.0"
        for k in $( eval echo {1..$numberOfRuns} )
        do
            rm -rf logs/*
            # create population of only agent
            # (could be changed in future for other groups agents + this agent)
            make run
            logDir=("logs/*")
            lifeFile=$logDir"/main.json"
            deathFile=$logDir"/death.json"
            # pass in logfile, num agents, agent_config.json, bestAgent.config, current iteration to python script
            lifeExpectancy=$(python3 pkg/agents/team4/agent1/getLifeExpectancy.py $lifeFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)
            deathRate=$(python3 pkg/agents/team4/agent1/getDeathRate.py $deathFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)
            # record average survival rate
            # echo "      Run " $k " - " $pythonOutput 
            # individualResults="$individualResults + $pythonOutput" | python3
            averageLifeExpectancy=`echo $averageLifeExpectancy+$lifeExpectancy | bc`
            averageDeathRate=`echo $averageDeathRate+$deathRate | bc`
        done
        averageLifeExpectancy=`echo $averageLifeExpectancy/$numberOfRuns | bc -l` 
        averageDeathRate=`echo $averageDeathRate/$numberOfRuns | bc -l` 
        # echo "  Average perfomance " $averageLifeExpectancy 
        arrLifeExp+=($averageLifeExpectancy)
        arrDeathRate+=($averageDeathRate)
    done
    printf -v joinedLifeExp '%s,' ${arrLifeExp[*]}
    echo "[${joinedLifeExp%,}]" > $agentLifeExpectanciesFile

    printf -v joinedDeathRate '%s,' ${arrDeathRate[*]}
    echo "[${joinedDeathRate%,}]" > $agentDeathRateFile

    # generate new set of best agents generated from previous perfomance 
    python3 pkg/agents/team4/agent1/generateNewBestAgents.py $bestAgentsFile $agentLifeExpectanciesFile $numberOfHealthLevels $agentDeathRateFile
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