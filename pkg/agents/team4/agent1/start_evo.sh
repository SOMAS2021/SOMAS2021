#!/bin/bash

selfish='false'

print_usage() {
  printf "Usage: ... \n"
  printf "s) set to train for selfish agent \n"
  printf "h) get help \n"
}

while getopts 'sh' flag
do
  case "${flag}" in
    s) selfish='true' ;;
    h) print_usage
       exit 1 ;;
  esac
done

# Navigating to root directory
rootdirpath="../../../../"
cd $rootdirpath

# Initialising config files
agentConfigFile="pkg/agents/team4/agent1/agentConfig.json"
bestAgentsFile="pkg/agents/team4/agent1/bestAgents.json"
agentLifeExpectanciesFile="pkg/agents/team4/agent1/agentLifeExpectancies.json"
agentOurLifeExpectanciesFile="pkg/agents/team4/agent1/agentOurLifeExpectancies.json"
agentDeathRateFile="pkg/agents/team4/agent1/agentDeathRate.json"
rm $agentConfigFile $bestAgentsFile
touch $agentConfigFile $bestAgentsFile

numberOfHealthLevels=4
numberOfBestAgents=5
numberOfAgentsPerSim=5
numberOfIterations=5
numberOfRuns=1

# Generate set of agents with 0 parameters
python3 pkg/agents/team4/agent1/initaliseConfig.py $agentConfigFile $bestAgentsFile $numberOfHealthLevels $numberOfBestAgents

for i in $( eval echo {1..$numberOfIterations} )
do
    echo "ITERATION " $i
    echo ""
    arrLifeExp=()
    arrOurLifeExp=()
    arrDeathRate=()
    for j in $( eval echo {1..$numberOfBestAgents} )
    do
        echo "  Getting average performance of agent " $j
        averageLifeExpectancy="0.0"
        averageOurLifeExpectancy="0.0"
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

            # Set space as the delimiter
            OLDIFS=$IFS
            IFS=';'
            #Read the split words into an array based on space delimiter
            read -a agentLifeExpectanciesArray <<< "$lifeExpectancy"
            IFS=$OLDIFS
            
            deathRate=$(python3 pkg/agents/team4/agent1/getDeathRate.py $deathFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)

            averageLifeExpectancy=`echo $averageLifeExpectancy+${agentLifeExpectanciesArray[0]} | bc`

            averageOurLifeExpectancy=`echo $averageOurLifeExpectancy+${agentLifeExpectanciesArray[1]} | bc`

            averageDeathRate=`echo $averageDeathRate+$deathRate | bc`

        done
        averageLifeExpectancy=`echo $averageLifeExpectancy/$numberOfRuns | bc -l`
        averageOurLifeExpectancy=`echo $averageOurLifeExpectancy/$numberOfRuns | bc -l` 
        averageDeathRate=`echo $averageDeathRate/$numberOfRuns | bc -l` 
        arrLifeExp+=($averageLifeExpectancy)
        arrOurLifeExp+=($averageOurLifeExpectancy)
        arrDeathRate+=($averageDeathRate)
    done
    printf -v joinedLifeExp '%s,' ${arrLifeExp[*]}
    echo "[${joinedLifeExp%,}]" > $agentLifeExpectanciesFile

    printf -v joinedOurLifeExp '%s,' ${arrOurLifeExp[*]}
    echo "[${joinedOurLifeExp%,}]" > $agentOurLifeExpectanciesFile

    printf -v joinedDeathRate '%s,' ${arrDeathRate[*]}
    echo "[${joinedDeathRate%,}]" > $agentDeathRateFile

    # generate new set of best agents generated from previous perfomance 
    python3 pkg/agents/team4/agent1/generateNewBestAgents.py $bestAgentsFile $agentLifeExpectanciesFile $numberOfHealthLevels $agentDeathRateFile $agentOurLifeExpectanciesFile $selfish
    echo "------------------------------------------"
done

# TODO: autmatate finding best agents (writing to a different file based on the selfish flag)
# python3 findBestConfigs.py

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