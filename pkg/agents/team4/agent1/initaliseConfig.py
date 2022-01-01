import json
import sys
import random

agent_config_file_name = sys.argv[1]
best_agent_file_name = sys.argv[2]
number_of_health_levels = sys.argv[3]
number_of_best_agents = sys.argv[4]

agent_config_file = open(agent_config_file_name, 'r+')
best_agents_file = open(best_agent_file_name, 'r+')

number_of_agents = int(number_of_best_agents) + 1

# Initialising all configuration files to 0 values
best_agents = []
for _ in range(number_of_agents):
    agent_config = {
        "FoodToEat": [],
        "DaysToWait": [],
    }

    # random.randint(0, 100)
    agent_config["FoodToEat"] = [0.0
                                 for _ in range(int(number_of_health_levels))]

    # random.randint(0, 7)
    agent_config["DaysToWait"] = [0.0
                                  for _ in range(int(number_of_health_levels))]

    best_agents.append(agent_config)

agent_config_str = json.dumps(best_agents[0], indent=4)
agent_config_file.write(agent_config_str)
# print("Initialising agent config to: {0}".format(agent_config_str))

best_agents_str = json.dumps(best_agents, indent=4)
best_agents_file.write(best_agents_str)
# print("Initialising best agents to: {0}".format(best_agents_str))

agent_config_file.close()
best_agents_file.close()
