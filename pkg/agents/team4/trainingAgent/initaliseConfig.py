import json
import sys
import random

agent_config_file_name = sys.argv[1]
best_agent_file_name = sys.argv[2]
number_of_health_levels = sys.argv[3]
number_of_best_agents = sys.argv[4]

agent_config_file = open(agent_config_file_name, 'r+')
best_agents_file = open(best_agent_file_name, 'r+')

number_of_agents = int(number_of_best_agents)

# Initialising all configuration files
best_agents = []


for i in range(number_of_agents):
    # agent config template
    agent_config = {
        "FoodToEat": [],
        "WaitProbability": [],
    }

    # increase food to eat by increments of 20
    agent_config["FoodToEat"] = [(i % 5) * 20
                                 for _ in range(int(number_of_health_levels))]
    # set WaitProbability to a random number between 0 and 100
    agent_config["WaitProbability"] = [random.randint(0, 100)
                                       for _ in range(int(number_of_health_levels))]
    best_agents.append(agent_config)

for best_agent in best_agents:
    best_agent["FoodToEat"][0] = -1
    best_agent["WaitProbability"][0] = -1

agent_config_str = json.dumps(
    best_agents[0], indent=4)  # convert to JSON string
agent_config_file.write(agent_config_str)  # write to agentConfig.json file
# print("Initialising agent config to: {0}".format(agent_config_str))

best_agents_str = json.dumps(best_agents, indent=4)  # convert to JSON string
best_agents_file.write(best_agents_str)
# print("Initialising best agents to: {0}".format(best_agents_str))

agent_config_file.close()
best_agents_file.close()
