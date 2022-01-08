#!/bin/bash

# set up flags used for simulation config
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
agentConfigFile="pkg/agents/team4/trainingAgent/configs/agentConfig.json"
bestAgentsFile="pkg/agents/team4/trainingAgent/configs/bestAgents.json"
agentLifeExpectanciesFile="pkg/agents/team4/trainingAgent/configs/agentLifeExpectancies.json"
agentOurLifeExpectanciesFile="pkg/agents/team4/trainingAgent/configs/agentOurLifeExpectancies.json"
agentOtherLifeExpectanciesFile="pkg/agents/team4/trainingAgent/configs/agentOtherLifeExpectancies.json"
agentDeathRateFile="pkg/agents/team4/trainingAgent/configs/agentDeathRate.json"

#removing old files and creating new files
rm -rf "pkg/agents/team4/trainingAgent/configs"
mkdir "pkg/agents/team4/trainingAgent/configs"
touch $agentConfigFile $bestAgentsFile
mkdir "pkg/agents/team4/trainingAgent/configs/storedConfigs"

numberOfHealthLevels=4 # health bands based on HP
numberOfBestAgents=5 # number of agents in training population
numberOfAgentsPerSim=16 # total number of agents in the simulation
numberOfIterations=3 # number of training iterations
numberOfRuns=2 # number of runs to average over per iteration

# Initialise parameters of initial agent polulation
python3 pkg/agents/team4/trainingAgent/initaliseConfig.py $agentConfigFile $bestAgentsFile $numberOfHealthLevels $numberOfBestAgents 

for i in $( eval echo {1..$numberOfIterations} )
do
    echo "ITERATION " $i
    echo ""
    # create array for each stored metric
    arrLifeExp=()
    arrOurLifeExp=()
    arrOtherLifeExp=()
    arrDeathRate=()
    
    for j in $( eval echo {1..$numberOfBestAgents} )
    do
        echo "  Getting average performance of agent " $j
        
        # intialise average metrics
        averageLifeExpectancy="0.0"
        averageOtherLifeExpectancy="0.0"
        averageOurLifeExpectancy="0.0"
        averageDeathRate="0.0"
        
        for k in $( eval echo {1..$numberOfRuns} )
        do
            rm -rf logs/*
            # run particular agent from population with all other teams' agents
            make run
            logDir=("logs/*")
            lifeFile=$logDir"/main.json"
            deathFile=$logDir"/death.json"
            # pass in logfile, num agents, agent_config.json, bestAgent.config, current iteration to python script
            lifeExpectancy=$(python3 pkg/agents/team4/trainingAgent/getLifeExpectancy.py $lifeFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j)

            # handle python output to get life expectencies
            OLDIFS=$IFS
            IFS=';'
            read -a agentLifeExpectanciesArray <<< "$lifeExpectancy"
            IFS=$OLDIFS
            
            deathRate=$(python3 pkg/agents/team4/trainingAgent/getDeathRate.py $deathFile $numberOfAgentsPerSim $agentConfigFile $bestAgentsFile $j) # ------------- NOT USED YET ---------------

            # get running total for life expectancies to get the average
            averageLifeExpectancy=`echo $averageLifeExpectancy+${agentLifeExpectanciesArray[0]} | bc`
            averageOurLifeExpectancy=`echo $averageOurLifeExpectancy+${agentLifeExpectanciesArray[1]} | bc`
            averageOtherLifeExpectancy=`echo $averageOtherLifeExpectancy+${agentLifeExpectanciesArray[2]} | bc`

            averageDeathRate=`echo $averageDeathRate+$deathRate | bc` # ------------- NOT USED YET ---------------
        done
        averageLifeExpectancy=`echo $averageLifeExpectancy/$numberOfRuns | bc -l`
        averageOurLifeExpectancy=`echo $averageOurLifeExpectancy/$numberOfRuns | bc -l` 
        averageOtherLifeExpectancy=`echo $averageOtherLifeExpectancy/$numberOfRuns | bc -l`
        averageDeathRate=`echo $averageDeathRate/$numberOfRuns | bc -l` # ------------- NOT USED YET ---------------


        arrLifeExp+=($averageLifeExpectancy)
        arrOurLifeExp+=($averageOurLifeExpectancy)
        arrDeathRate+=($averageDeathRate) # ------------- NOT USED YET ---------------
        arrOtherLifeExp+=($averageOtherLifeExpectancy)
    done
    printf -v joinedLifeExp '%s,' ${arrLifeExp[*]}
    echo "[${joinedLifeExp%,}]" > $agentLifeExpectanciesFile

    printf -v joinedOurLifeExp '%s,' ${arrOurLifeExp[*]}
    echo "[${joinedOurLifeExp%,}]" > $agentOurLifeExpectanciesFile

    printf -v joinedOtherLifeExp '%s,' ${arrOtherLifeExp[*]}
    echo "[${joinedOtherLifeExp%,}]" > $agentOtherLifeExpectanciesFile

    # ------------- NOT USED YET ---------------
    printf -v joinedDeathRate '%s,' ${arrDeathRate[*]}
    echo "[${joinedDeathRate%,}]" > $agentDeathRateFile

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