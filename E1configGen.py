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