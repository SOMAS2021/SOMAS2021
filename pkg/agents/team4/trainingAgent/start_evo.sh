#!/bin/bash

selfish='False'
selfless='False'

print_usage() {
  printf "Usage: ... \n"
  printf "f) set to train for selfish agent \n"
  printf "l) set to train for selfless agent \n"
  printf "h) get help \n"
}

while getopts 'fhl' flag
do
  case "${flag}" in
    f) selfish='True' ;;
    l) selfless='True' ;;
    h) print_usage
       exit 1 ;;
  esac
done

# Navigating to root directory
rootdirpath="../../../../"
cd $rootdirpath

# Initialising config files
agentConfigFile="pkg/agents/team4/trainingAgent/agentConfig.json"
bestAgentsFile="pkg/agents/team4/trainingAgent/bestAgents.json"
agentLifeExpectanciesFile="pkg/agents/team4/trainingAgent/agentLifeExpectancies.json"
agentOurLifeExpectanciesFile="pkg/agents/team4/trainingAgent/agentOurLifeExpectancies.json"
agentOtherLifeExpectanciesFile="pkg/agents/team4/trainingAgent/agentOtherLifeExpectancies.json"
agentDeathRateFile="pkg/agents/team4/trainingAgent/agentDeathRate.json"
rm $agentConfigFile $bestAgentsFile
touch $agentConfigFile $bestAgentsFile
rm -rf "pkg/agents/team4/trainingAgent/storedagents"
mkdir "pkg/agents/team4/trainingAgent/storedagents"

numberOfHealthLevels=4
numberOfBestAgents=5
numberOfAgentsPerSim=16
numberOfIterations=15
numberOfRuns=1

# Generate set of agents with 0 parameters
python3 pkg/agents/team4/trainingAgent/initaliseConfig.py $agentConfigFile $bestAgentsFile $numberOfHealthLevels $numberOfBestAgents

for i in $( eval echo {1..$numberOfIterations} )
do
    echo "ITERATION " $i
    echo ""
    arrLifeExp=()
    arrOurLifeExp=()
    arrOtherLifeExp=()
    arrDeathRate=()
    for j in $( eval echo {1..$numberOfBestAgents} )
    do
        echo "  Getting average performance of agent " $j
        averageLifeExpectancy="0.0"
        averageOtherLifeExpectancy="0.0"
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
            lifeExpectancy=$(python3 pkg/agents/team4/trainingAgent/getLifeExpectancy.py $lifeFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)

            # Set space as the delimiter
            OLDIFS=$IFS
            IFS=';'
            #Read the split words into an array based on space delimiter
            read -a agentLifeExpectanciesArray <<< "$lifeExpectancy"
            IFS=$OLDIFS
            
            deathRate=$(python3 pkg/agents/team4/trainingAgent/getDeathRate.py $deathFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)

            averageLifeExpectancy=`echo $averageLifeExpectancy+${agentLifeExpectanciesArray[0]} | bc`

            averageOurLifeExpectancy=`echo $averageOurLifeExpectancy+${agentLifeExpectanciesArray[1]} | bc`

            averageOtherLifeExpectancy=`echo $averageOtherLifeExpectancy+${agentLifeExpectanciesArray[2]} | bc`

            averageDeathRate=`echo $averageDeathRate+$deathRate | bc`

        done
        averageLifeExpectancy=`echo $averageLifeExpectancy/$numberOfRuns | bc -l`
        averageOurLifeExpectancy=`echo $averageOurLifeExpectancy/$numberOfRuns | bc -l` 
        averageDeathRate=`echo $averageDeathRate/$numberOfRuns | bc -l` 
        averageOtherLifeExpectancy=`echo $averageOtherLifeExpectancy/$numberOfRuns | bc -l`


        arrLifeExp+=($averageLifeExpectancy)
        arrOurLifeExp+=($averageOurLifeExpectancy)
        arrDeathRate+=($averageDeathRate)
        arrOtherLifeExp+=($averageOtherLifeExpectancy)
    done
    printf -v joinedLifeExp '%s,' ${arrLifeExp[*]}
    echo "[${joinedLifeExp%,}]" > $agentLifeExpectanciesFile

    printf -v joinedOurLifeExp '%s,' ${arrOurLifeExp[*]}
    echo "[${joinedOurLifeExp%,}]" > $agentOurLifeExpectanciesFile

    printf -v joinedDeathRate '%s,' ${arrDeathRate[*]}
    echo "[${joinedDeathRate%,}]" > $agentDeathRateFile

    printf -v joinedOtherLifeExp '%s,' ${arrOtherLifeExp[*]}
    echo "[${joinedOtherLifeExp%,}]" > $agentOtherLifeExpectanciesFile

    # generate new set of best agents generated from previous perfomance 
    python3 pkg/agents/team4/trainingAgent/generateNewBestAgents.py $bestAgentsFile $agentLifeExpectanciesFile $agentDeathRateFile $agentOurLifeExpectanciesFile $agentOtherLifeExpectanciesFile $numberOfHealthLevels $selfish $selfless $i
    echo "------------------------------------------"
done

# TODO: autmatate finding best agents (writing to a different file based on the selfish flag)
# python3 findBestConfigs.py

# TODO: save lowest floor that we've ever been on and use that to optimise the food to eat and days to wait (writing to a different file based on the selfish flag)

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