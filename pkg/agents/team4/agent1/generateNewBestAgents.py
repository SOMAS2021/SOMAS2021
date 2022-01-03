import sys
import random
import json
import numpy as np

best_agents_file_name = sys.argv[1]
agent_life_expectencies_file_name = sys.argv[2]
number_of_health_levels = int(sys.argv[3])
agent_death_rate_file_name = sys.argv[4]

# read bestAgent.config
best_agents_file = open(best_agents_file_name)
best_agents_list = json.load(best_agents_file)
best_agents_file.close

# read the array of average time survied by each agent
agent_life_expectencies_file = open(agent_life_expectencies_file_name)
performance_list_life = json.load(agent_life_expectencies_file)
agent_life_expectencies_file.close()

# read the array of average death by each agent
agent_death_rate_file = open(agent_death_rate_file_name)
performance_list_death = json.load(agent_death_rate_file)
agent_death_rate_file.close()

# sort array but keep a track of index to obtain best agent
perfomance_list_indices = range(len(performance_list_life))
zipped_performance = zip(performance_list_life, performance_list_death, perfomance_list_indices)
sorted_zipped_performance = sorted(zipped_performance, reverse=True)
tuples = zip(*sorted_zipped_performance)
sorted_performance_list_life, sorted_performance_list_death, sorted_performance_list_indices = [
    list(tuple) for tuple in tuples]

print("The current life expectencies are: {0}".format(
    sorted_performance_list_life))

# get the 3 best agents
agent_1 = best_agents_list[sorted_performance_list_indices[0]]
agent_2 = best_agents_list[sorted_performance_list_indices[1]]
agent_3 = best_agents_list[sorted_performance_list_indices[2]]
agent_4 = best_agents_list[sorted_performance_list_indices[3]]
agent_5 = best_agents_list[sorted_performance_list_indices[4]]

################################################################################################
#life[0] death[0] timestamp
agents_of_time = [agent_1,agent_2,agent_3,agent_4,agent_5]
outfile = []
for i in range(len(agents_of_time)):
    agents_of_time[i]['life'] = performance_list_life[sorted_performance_list_indices[i]]
    agents_of_time[i]['death'] = performance_list_death[sorted_performance_list_indices[i]]

from datetime import datetime
timestamp = datetime.now(tz=None)
file_name_out= "pkg/agents/team4/agent1/"+"{0}/{1}_{2}_{3}.json".format("storedagents", timestamp, performance_list_life[0],performance_list_death[0])
file_name_out=file_name_out.replace(" ","_")
old_agents_file = open(file_name_out, 'w')
old_agents_file.write(json.dumps(agents_of_time, indent=4))
old_agents_file.close()
################################################################################################


attribute = ["FoodToEat", "DaysToWait"]
noise_multiplier_for_attribute = {
    "FoodToEat": 5,
    "DaysToWait": 2
}

mu, sigma = 0, 1  # mean and standard deviation

# create mutated agents based on the best 3 agents
agent_1_slight_mutations_1 = best_agents_list[sorted_performance_list_indices[0]]
for i in attribute:
    for j in range(number_of_health_levels):
        rand = int(
            noise_multiplier_for_attribute[i]*np.random.normal(mu, sigma))
        add_sub = random.random()
        if (add_sub < 0.5):
            agent_1_slight_mutations_1[i][j] += rand
        else:
            agent_1_slight_mutations_1[i][j] -= rand

        if (attribute == "FoodToEat" and agent_1_slight_mutations_1[i][j] > 100):
            agent_1_slight_mutations_1[i][j] = 100
        if (agent_1_slight_mutations_1[i][j] < 0):
            agent_1_slight_mutations_1[i][j] = 0

# UNUSED FOR NOW
# agent_1_slight_mutations_2 = agent_1
# for i in equation:
#     for j in range(number_of_health_levels):
#         rand = 0.001*random.random()
#         add_sub = random.random()
#         if (add_sub < 0.5):
#             agent_1_slight_mutations_2[i][j] += rand
#         else:
#             agent_1_slight_mutations_2[i][j] -= rand

#         if (agent_1_slight_mutations_2[i][j] > 1):
#             agent_1_slight_mutations_2[i][j] = 1
#         elif(agent_1_slight_mutations_2[i][j] < -1):
#             agent_1_slight_mutations_2[i][j] = -1

# agent_1_big_mutation = agent_1
# for i in equation:
#     for j in range(number_of_health_levels):
#         rand = 0.1*random.random()
#         add_sub = random.random()
#         if (add_sub < 0.5):
#             agent_1_big_mutation[i][j] += rand
#         else:
#             agent_1_big_mutation[i][j] -= rand

#         if (agent_1_big_mutation[i][j] > 1):
#             agent_1_big_mutation[i][j] = 1
#         elif(agent_1_big_mutation[i][j] < -1):
#             agent_1_big_mutation[i][j] = -1

agent_mix_1 = {}
for i in attribute:
    values = []
    for j in range(number_of_health_levels):
        temp_val = 0
        rand = random.random()
        if (rand < 0.66):
            temp_val = agent_1[i][j]
        elif (rand < 0.90):
            temp_val = agent_2[i][j]
        else:
            temp_val = agent_3[i][j]
        values.append(temp_val)

    agent_mix_1[i] = values

agent_mix_2 = {}
for i in attribute:
    values = []
    for j in range(number_of_health_levels):
        temp_val = 0.0
        rand = random.random()
        if (rand < 0.66):
            temp_val = agent_1[i][j]
        elif (rand < 0.90):
            temp_val = agent_2[i][j]
        else:
            temp_val = agent_3[i][j]
        values.append(temp_val)

    agent_mix_2[i] = values

# agent_mix_3 = {}
# for i in attribute:
#     values = []
#     for j in range(number_of_health_levels):
#         temp_val = 0.0
#         rand = random.random()
#         if (rand < 0.33):
#             temp_val = agent_1[i][j]
#         elif (rand < 0.66):
#             temp_val = agent_2[i][j]
#         else:
#             temp_val = agent_3[i][j]
#         values.append(temp_val)

#     agent_mix_3[i] = values

random_mutation_agent = {
    "FoodToEat": [],
    "DaysToWait": [],
}

random_mutation_agent["FoodToEat"] = [random.randint(0, 100)
                                      for _ in range(int(number_of_health_levels))]

random_mutation_agent["DaysToWait"] = [random.randint(0, 7)
                                       for _ in range(int(number_of_health_levels))]


new_best_agents = [agent_1, agent_1_slight_mutations_1, agent_mix_1,
                   agent_mix_2, random_mutation_agent]

for best_agent in new_best_agents:
    best_agent["FoodToEat"][0] = -1


# #life[0] death[0] timestamp
# agents_of_time = [agent_1,agent_2,agent_3,agent_4,agent_5]
# outfile = []
# for i in range(len(agents_of_time)):
#     agents_of_time[i]['life'] = performance_list_life[i]
#     agents_of_time[i]['death'] = performance_list_death[i]

# from datetime import datetime
# timestamp = datetime.now(tz=None)
# file_name_out= "pkg/agents/team4/agent1/"+"{0}/{1}_{2}_{3}.json".format("storedagents", timestamp, performance_list_life[0],performance_list_death[0])
# file_name_out=file_name_out.replace(" ","_")
# old_agents_file = open(file_name_out, 'w')
# old_agents_file.write(json.dumps(agents_of_time, indent=4))
# old_agents_file.close()

best_agents_file = open(best_agents_file_name, 'w')
best_agents_file.write(json.dumps(new_best_agents, indent=4))
best_agents_file.close()
