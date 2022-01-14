UseFoodPerAgentRatio = "false" 
FoodOnPlatform = 100 
FoodPerAgentRatio = 0 
Team1Agents = 2 
Team1Agents2 = 0
Team2Agents = 2 
Team3Agents = 2 
Team4Agents = 2 
Team5Agents = 2 
Team6Agents = 2 
Team7Agents = 2 
RandomAgents = 2 
AgentHP = 100 
AgentsPerFloor = 1 
TicksPerFloor = 10 
SimDays = 8 
ReshuffleDays = 7
maxHP = 100 
weakLevel = 10 
width = 48 
tau = 15 
MaxFoodIntake = 60 
hpReqCToW = 5 
hpCritical = 3 
maxDayCritical = 3 
HPLossBase = 5 
HPLossSlope = 0.2 
LogMain = "false" 
LogStory = "false" 
SimTimeoutSeconds = 60

#Experiment 1
print("Generating Experiment 1 Configs")
for foodPerAgent in [5]:
    for SelfishNum in [0, 3]:
        for RandomNum in [0, 3]:
            for OtherNum in [3]:
                UseFoodPerAgentRatio = "true"
                FoodOnPlatform = 100 
                FoodPerAgentRatio = foodPerAgent
                Team1Agents = SelfishNum
                Team2Agents = OtherNum
                Team3Agents = OtherNum
                Team4Agents = OtherNum
                Team5Agents = OtherNum
                Team6Agents = OtherNum
                Team7Agents = OtherNum
                RandomAgents = RandomNum
                TicksPerFloor = 10 
                AgentHP = 100 
                SimDays = 400
                ReshuffleDays = 7
                SimTimeoutSeconds = 1000

                SimName = f"E1-{SelfishNum}-{RandomNum}-{OtherNum}-{FoodPerAgentRatio}"
                logName = f"log-{SimName}"
                configName = f"config-{SimName}"
                jsonString = f"""
                    "UseFoodPerAgentRatio": {UseFoodPerAgentRatio},
                    "FoodOnPlatform": {FoodOnPlatform},
                    "FoodPerAgentRatio": {FoodPerAgentRatio},
                    "Team1Agents": {Team1Agents},
                    "Team1Agents2": {Team1Agents2},
                    "Team2Agents": {Team2Agents},
                    "Team3Agents": {Team3Agents},
                    "Team4Agents": {Team4Agents},
                    "Team5Agents": {Team5Agents},
                    "Team6Agents": {Team6Agents},
                    "Team7Agents": {Team7Agents},
                    "RandomAgents": {RandomAgents},
                    "AgentHP": {AgentHP},
                    "AgentsPerFloor": {AgentsPerFloor},
                    "TicksPerFloor": {TicksPerFloor},
                    "SimDays": {SimDays},
                    "ReshuffleDays": {ReshuffleDays},
                    "maxHP": {maxHP},
                    "weakLevel": {weakLevel},
                    "width": {width},
                    "tau": {tau},
                    "MaxFoodIntake": {MaxFoodIntake},
                    "hpReqCToW": {hpReqCToW},
                    "hpCritical": {hpCritical},
                    "maxDayCritical": {maxDayCritical},
                    "HPLossBase": {HPLossBase},
                    "HPLossSlope": {HPLossSlope},
                    "LogMain": {LogMain},
                    "LogStory": {LogStory},
                    "SimTimeoutSeconds": {SimTimeoutSeconds},
                    "LogFileName": {logName}
                """
                jsonString = "{" + jsonString + "}"
                f = open(f"configs/{configName}.json", "w")
                f.write(jsonString)
                f.close()

# Experiment 2
print("Generating Experiment 2 Configs")
for reshuffleFrequency in [2,7,14]: # effect of reshuffleFrequency
    for platformTicksPerFloor in [1,5,10]: # limiting actions
        for foodPerAgent in [5]:
            for SelfishNum in [0]:
                for RandomNum in [0]:
                    for OtherNum in [3]:
                        UseFoodPerAgentRatio = "true"
                        FoodOnPlatform = 100 
                        FoodPerAgentRatio = foodPerAgent
                        Team1Agents = SelfishNum
                        Team2Agents = OtherNum
                        Team3Agents = OtherNum
                        Team4Agents = OtherNum
                        Team5Agents = OtherNum
                        Team6Agents = OtherNum
                        Team7Agents = OtherNum
                        RandomAgents = RandomNum
                        TicksPerFloor = platformTicksPerFloor 
                        AgentHP = 100 
                        SimDays = 400
                        ReshuffleDays = reshuffleFrequency
                        SimTimeoutSeconds = 1000

                        SimName = f"E2-{platformTicksPerFloor}-{reshuffleFrequency}"
                        logName = f"log-{SimName}"
                        configName = f"config-{SimName}"
                        
                        jsonString = f"""
                            "UseFoodPerAgentRatio": {UseFoodPerAgentRatio},
                            "FoodOnPlatform": {FoodOnPlatform},
                            "FoodPerAgentRatio": {FoodPerAgentRatio},
                            "Team1Agents": {Team1Agents},
                            "Team1Agents2": {Team1Agents2},
                            "Team2Agents": {Team2Agents},
                            "Team3Agents": {Team3Agents},
                            "Team4Agents": {Team4Agents},
                            "Team5Agents": {Team5Agents},
                            "Team6Agents": {Team6Agents},
                            "Team7Agents": {Team7Agents},
                            "RandomAgents": {RandomAgents},
                            "AgentHP": {AgentHP},
                            "AgentsPerFloor": {AgentsPerFloor},
                            "TicksPerFloor": {TicksPerFloor},
                            "SimDays": {SimDays},
                            "ReshuffleDays": {ReshuffleDays},
                            "maxHP": {maxHP},
                            "weakLevel": {weakLevel},
                            "width": {width},
                            "tau": {tau},
                            "MaxFoodIntake": {MaxFoodIntake},
                            "hpReqCToW": {hpReqCToW},
                            "hpCritical": {hpCritical},
                            "maxDayCritical": {maxDayCritical},
                            "HPLossBase": {HPLossBase},
                            "HPLossSlope": {HPLossSlope},
                            "LogMain": {LogMain},
                            "LogStory": {LogStory},
                            "SimTimeoutSeconds": {SimTimeoutSeconds},
                            "LogFileName": {logName}
                        """
                        jsonString = "{" + jsonString + "}"
                        f = open(f"configs/{configName}.json", "w")
                        f.write(jsonString)
                        f.close()

#Experiemnt 3
print("Generating Experiment 3 Configs")
Tournament_array = [0,0,0,0,0,0,0,0]
# Generate Tournament Configs
for foodPerAgent in [5]:
    for playerA in [0,1,2,3,4,5,6,7]:
        possibleOpponentsArray = [0,1,2,3,4,5,6,7]
        possibleOpponentsArray = possibleOpponentsArray[playerA:8]
        
        for playerB in possibleOpponentsArray:
            
            if(playerA != playerB):
                
                Tournament_array = [0,0,0,0,0,0,0,0]
                Tournament_array[playerA] = 10
                Tournament_array[playerB] = 10
                
                UseFoodPerAgentRatio = "true"
                FoodOnPlatform = 100 
                FoodPerAgentRatio = foodPerAgent
                RandomAgents = Tournament_array[0]
                Team1Agents  = Tournament_array[1]
                Team2Agents  = Tournament_array[2]
                Team3Agents  = Tournament_array[3]
                Team4Agents  = Tournament_array[4]
                Team5Agents  = Tournament_array[5]
                Team6Agents  = Tournament_array[6]
                Team7Agents  = Tournament_array[7]
                
                TicksPerFloor = 10 
                AgentHP = 100 
                SimDays = 400
                ReshuffleDays = 7
                SimTimeoutSeconds = 1000

                SimName = f"E3-{playerA}-{playerB}"
                logName = f"log-{SimName}"
                configName = f"config-{SimName}"

                jsonString = f"""
                    "UseFoodPerAgentRatio": {UseFoodPerAgentRatio},
                    "FoodOnPlatform": {FoodOnPlatform},
                    "FoodPerAgentRatio": {FoodPerAgentRatio},
                    "Team1Agents": {Team1Agents},
                    "Team1Agents2": {Team1Agents2},
                    "Team2Agents": {Team2Agents},
                    "Team3Agents": {Team3Agents},
                    "Team4Agents": {Team4Agents},
                    "Team5Agents": {Team5Agents},
                    "Team6Agents": {Team6Agents},
                    "Team7Agents": {Team7Agents},
                    "RandomAgents": {RandomAgents},
                    "AgentHP": {AgentHP},
                    "AgentsPerFloor": {AgentsPerFloor},
                    "TicksPerFloor": {TicksPerFloor},
                    "SimDays": {SimDays},
                    "ReshuffleDays": {ReshuffleDays},
                    "maxHP": {maxHP},
                    "weakLevel": {weakLevel},
                    "width": {width},
                    "tau": {tau},
                    "MaxFoodIntake": {MaxFoodIntake},
                    "hpReqCToW": {hpReqCToW},
                    "hpCritical": {hpCritical},
                    "maxDayCritical": {maxDayCritical},
                    "HPLossBase": {HPLossBase},
                    "HPLossSlope": {HPLossSlope},
                    "LogMain": {LogMain},
                    "LogStory": {LogStory},
<<<<<<< HEAD
                    "SimTimeoutSeconds": {SimTimeoutSeconds},
                    "LogFileName": {logName}
                """
                
                jsonString = "{" + jsonString + "}"
                f = open(f"configs/{configName}.json", "w")
                f.write(jsonString)
                f.close()

# Experiment 4
print("Generating Experiment 4 Configs")
for playerA in [2,3,4,5,6,7]:
    for numSelfish in [0,1,2,3,4,5,6,7,8,9,10]:
        
        Tournament_array = [0,numSelfish,0,0,0,0,0,0]
        Tournament_array[playerA] = 10
        
        UseFoodPerAgentRatio = "true"
        FoodOnPlatform = 100 
        FoodPerAgentRatio = foodPerAgent
        RandomAgents = Tournament_array[0]
        Team1Agents  = Tournament_array[1]
        Team2Agents  = Tournament_array[2]
        Team3Agents  = Tournament_array[3]
        Team4Agents  = Tournament_array[4]
        Team5Agents  = Tournament_array[5]
        Team6Agents  = Tournament_array[6]
        Team7Agents  = Tournament_array[7]
        
        TicksPerFloor = 10 
        AgentHP = 100 
        SimDays = 400
        ReshuffleDays = 7
        SimTimeoutSeconds = 1000

        SimName = f"E4-selfish-{playerA}-{numSelfish}"
        logName = f"log-{SimName}"
        configName = f"config-{SimName}"

        jsonString = f"""
            "UseFoodPerAgentRatio": {UseFoodPerAgentRatio},
            "FoodOnPlatform": {FoodOnPlatform},
            "FoodPerAgentRatio": {FoodPerAgentRatio},
            "Team1Agents": {Team1Agents},
            "Team1Agents2": {Team1Agents2},
            "Team2Agents": {Team2Agents},
            "Team3Agents": {Team3Agents},
            "Team4Agents": {Team4Agents},
            "Team5Agents": {Team5Agents},
            "Team6Agents": {Team6Agents},
            "Team7Agents": {Team7Agents},
            "RandomAgents": {RandomAgents},
            "AgentHP": {AgentHP},
            "AgentsPerFloor": {AgentsPerFloor},
            "TicksPerFloor": {TicksPerFloor},
            "SimDays": {SimDays},
            "ReshuffleDays": {ReshuffleDays},
            "maxHP": {maxHP},
            "weakLevel": {weakLevel},
            "width": {width},
            "tau": {tau},
            "MaxFoodIntake": {MaxFoodIntake},
            "hpReqCToW": {hpReqCToW},
            "hpCritical": {hpCritical},
            "maxDayCritical": {maxDayCritical},
            "HPLossBase": {HPLossBase},
            "HPLossSlope": {HPLossSlope},
            "LogMain": {LogMain},
            "LogStory": {LogStory},
            "SimTimeoutSeconds": {SimTimeoutSeconds},
            "LogFileName": {logName}
        """
        
        jsonString = "{" + jsonString + "}"
        f = open(f"configs/{configName}.json", "w")
        f.write(jsonString)
        f.close()
playerA = 0
for playerA in [2,3,4,5,6,7]:
    for numRand in [0,1,2,3,4,5,6,7,8,9,10]:
        
        Tournament_array = [numRand,0,0,0,0,0,0,0]
        Tournament_array[playerA] = 10
        
        UseFoodPerAgentRatio = "true"
        FoodOnPlatform = 100 
        FoodPerAgentRatio = foodPerAgent
        RandomAgents = Tournament_array[0]
        Team1Agents  = Tournament_array[1]
        Team2Agents  = Tournament_array[2]
        Team3Agents  = Tournament_array[3]
        Team4Agents  = Tournament_array[4]
        Team5Agents  = Tournament_array[5]
        Team6Agents  = Tournament_array[6]
        Team7Agents  = Tournament_array[7]
        
        TicksPerFloor = 10 
        AgentHP = 100 
        SimDays = 400
        ReshuffleDays = 7
        SimTimeoutSeconds = 1000

        SimName = f"E4-random-{playerB}-{numRand}"
        logName = f"log-{SimName}"
        configName = f"config-{SimName}"

        jsonString = f"""
            "UseFoodPerAgentRatio": {UseFoodPerAgentRatio},
            "FoodOnPlatform": {FoodOnPlatform},
            "FoodPerAgentRatio": {FoodPerAgentRatio},
            "Team1Agents": {Team1Agents},
            "Team1Agents2": {Team1Agents2},
            "Team2Agents": {Team2Agents},
            "Team3Agents": {Team3Agents},
            "Team4Agents": {Team4Agents},
            "Team5Agents": {Team5Agents},
            "Team6Agents": {Team6Agents},
            "Team7Agents": {Team7Agents},
            "RandomAgents": {RandomAgents},
            "AgentHP": {AgentHP},
            "AgentsPerFloor": {AgentsPerFloor},
            "TicksPerFloor": {TicksPerFloor},
            "SimDays": {SimDays},
            "ReshuffleDays": {ReshuffleDays},
            "maxHP": {maxHP},
            "weakLevel": {weakLevel},
            "width": {width},
            "tau": {tau},
            "MaxFoodIntake": {MaxFoodIntake},
            "hpReqCToW": {hpReqCToW},
            "hpCritical": {hpCritical},
            "maxDayCritical": {maxDayCritical},
            "HPLossBase": {HPLossBase},
            "HPLossSlope": {HPLossSlope},
            "LogMain": {LogMain},
            "LogStory": {LogStory},
            "SimTimeoutSeconds": {SimTimeoutSeconds},
            "LogFileName": {logName}
        """
        
        jsonString = "{" + jsonString + "}"
        f = open(f"configs/{configName}.json", "w")
        f.write(jsonString)
<<<<<<< HEAD
        f.close()
=======
        f.close()
>>>>>>> adding config genrator scripts and a way to run all of them
=======
                    "SimTimeoutSeconds": {SimTimeoutSeconds}
                """
                
                jsonString = "{" + jsonString + "}"
                f = open(f"configs/configTournament-{playerA}-{playerB}-{FoodPerAgentRatio}.json", "w")
                f.write(jsonString)
                f.close()
>>>>>>> made come configs and added some scripts
