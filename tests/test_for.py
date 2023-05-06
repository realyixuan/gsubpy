sum = 0
for i in [1, 2, 3, 4]:
    sum += i

assert sum == 10


for i in range(4):
    if i == 2:
        break

assert i == 2


i = 0
total = 0
while i < 3:
    i += 1
    if i > 1:
        continue
    total += 1

assert total == 1

