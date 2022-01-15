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
                    "SimTimeoutSeconds": {SimTimeoutSeconds},
                    "LogFileName": {logName}
                """
                
                jsonString = "{" + jsonString + "}"
                f = open(f"configs/{configName}.json", "w")
                f.write(jsonString)
                f.close()
