import sys
import json

log_file_name = sys.argv[1]
number_of_agents = sys.argv[2]
agent_config_file_name = sys.argv[3]
best_agents_file_name = sys.argv[4]
current_iteration = sys.argv[5]

# reading contents of logfile
log_file = open(log_file_name, 'r')
lines = log_file.readlines()
log_file.close()

# convert log entries to json objects
logs = []
for i in range(len(lines)):
    logs.append(json.loads(lines[i]))

# ticks_per_floor = 10
# ticks_per_day = int(number_of_agents) * ticks_per_floor

# create dict of agents with all their corresponding log entries
agents = {}

# {"agentAge":33,"agentID":"506f71ff-0b92-4eb5-8fea-1750cb2d44f5","agentType":"Team5","level":"info","msg":"Agent survives till the end of the simulation","reporter":"simulation","time":"2022-01-03T11:17:44Z"}

for i in logs: # Add all agents to the dictionary, holds their age at point of deatrh or simulation end.
    try:
        if i['msg'] == "Killing agent":
            if i['agent_type'] not in agents:
                agents[i['agent_type']] = []
            agents[i['agent_type']].append(i["daysLived"])
        elif i['msg'] == "Agent survives till the end of the simulation":
            agents[i['agent_type']].append(i["agentAge"])
    except:
        continue

avgs = {}

# avg life expectancy of all agents globally
avg_of_all_agents = 0

# avg life expectancy of all agents not incl Team 4
avg_all_other_agents = 0

# number of agents not of Team 4
number_of_other_agents = 0

for agent in agents:
    # get avg life expectancy per agent and store
    avg = sum(agents[agent])/len(agents[agent])
    avgs[agent] = avg

    # increase global life exp
    avg_of_all_agents += avg

    if agent != "Team4":
        # count number of agents not Team 4
        number_of_other_agents += 1

        # increase global life exp of not team 4 
        avg_all_other_agents += avg


avg_of_all_agents /= len(avgs)

avg_all_other_agents /= number_of_other_agents

avg_days_lived = avgs["Team4"]

print(str(avg_of_all_agents)+";"+str(avg_days_lived)+";"+str(avg_all_other_agents))

# read best_agent file
best_agent_json_file = open(best_agents_file_name)
best_agent_json = json.load(best_agent_json_file)
best_agent_json_file.close()

try:
    # read agent at current_iteration+1
    next_agent = best_agent_json[int(current_iteration) + 1]
    # print("Changing agent config to: {0}".format(next_agent))

    # pass agent to agent_config file to create population for next run
    agent_config_file = open(agent_config_file_name, 'w')
    agent_config_file.write(json.dumps(next_agent, indent=4))
    agent_config_file.close()
except:
    pass

