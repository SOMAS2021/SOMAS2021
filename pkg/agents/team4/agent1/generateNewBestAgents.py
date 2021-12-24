import sys
import random
import json

best_agents_file_name = sys.argv[1]
agent_life_expectencies_file_name = sys.argv[2]
degree_of_equations = sys.argv[3]

# read bestAgent.config
best_agents_file = open(best_agents_file_name)
best_agents_list = json.load(best_agents_file)
best_agents_file.close

# read the array of average time survied by each agent
agent_life_expectencies_file = open(agent_life_expectencies_file_name)
performance_list = json.load(agent_life_expectencies_file)
agent_life_expectencies_file.close()

# sort array but keep a track of index to obtain best agent
perfomance_list_indices = range(len(performance_list))
zipped_performance = zip(performance_list, perfomance_list_indices)
sorted_zipped_performance = sorted(zipped_performance, reverse=True)
tuples = zip(*sorted_zipped_performance)
sorted_performance_list, sorted_performance_list_indices = [
    list(tuple) for tuple in tuples]

print("The current life expectencies are: {0}".format(
    sorted_performance_list))

# get the 3 best agents
agent_1 = best_agents_list[sorted_performance_list_indices[0]]
agent_2 = best_agents_list[sorted_performance_list_indices[1]]
agent_3 = best_agents_list[sorted_performance_list_indices[2]]

number_of_coefficients = int(degree_of_equations) + 1
equation = ["Floor", "Hp"]

# create mutated agents based on the best 3 agents
agent_1_slight_mutations_1 = agent_1
for i in equation:
    for j in range(number_of_coefficients):
        rand = 0.001*random.random()
        add_sub = random.random()
        if (add_sub < 0.5):
            agent_1_slight_mutations_1[i][j] += rand
        else:
            agent_1_slight_mutations_1[i][j] -= rand

        if (agent_1_slight_mutations_1[i][j] > 1):
            agent_1_slight_mutations_1[i][j] = 1
        elif(agent_1_slight_mutations_1[i][j] < -1):
            agent_1_slight_mutations_1[i][j] = -1

# UNUSED FOR NOW
# agent_1_slight_mutations_2 = agent_1
# for i in equation:
#     for j in range(number_of_coefficients):
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
#     for j in range(number_of_coefficients):
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
for i in equation:
    coefficients = []
    for j in range(number_of_coefficients):
        temp_coeff = 0.0
        rand = random.random()
        if (rand < 0.33):  # change this to 0.33
            temp_coeff = agent_1[i][j]
        elif (rand < 0.66):
            temp_coeff = agent_2[i][j]
        else:
            temp_coeff = agent_3[i][j]
        coefficients.append(temp_coeff)

    agent_mix_1[i] = coefficients

random_mutation_agent = {
    "Floor": [],
    "Hp": [],
}

agent_mix_2 = {}
for i in equation:
    coefficients = []
    for j in range(number_of_coefficients):
        temp_coeff = 0.0
        rand = random.random()
        if (rand < 0.33):
            temp_coeff = agent_1[i][j]
        elif (rand < 0.66):
            temp_coeff = agent_2[i][j]
        else:
            temp_coeff = agent_3[i][j]
        coefficients.append(temp_coeff)

    agent_mix_2[i] = coefficients

agent_mix_3 = {}
for i in equation:
    coefficients = []
    for j in range(number_of_coefficients):
        temp_coeff = 0.0
        rand = random.random()
        if (rand < 0.33):
            temp_coeff = agent_1[i][j]
        elif (rand < 0.66):
            temp_coeff = agent_2[i][j]
        else:
            temp_coeff = agent_3[i][j]
        coefficients.append(temp_coeff)

    agent_mix_3[i] = coefficients

random_mutation_agent = {
    "Floor": [],
    "Hp": [],
}

for k, _ in random_mutation_agent.items():
    random_mutation_agent[k] = [round(random.random(), 4)
                                for _ in range(number_of_coefficients)]

new_best_agents = [agent_1, agent_1_slight_mutations_1,
                   agent_mix_1, agent_mix_2, agent_mix_3]

best_agents_file = open(best_agents_file_name, 'w')
best_agents_file.write(json.dumps(new_best_agents, indent=4))
best_agents_file.close()
