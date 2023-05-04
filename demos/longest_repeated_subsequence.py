def LRS(X, m, n, lookup):
 
    if m == 0 or n == 0:
        return ''
 
    if X[m - 1] == X[n - 1] and m != n:
        return LRS(X, m - 1, n - 1, lookup) + X[m - 1]
    else:
        if lookup[m - 1][n] > lookup[m][n - 1]:
            return LRS(X, m - 1, n, lookup)
        else:
            return LRS(X, m, n - 1, lookup)
 
 
def LRSLength(X, lookup):
 
    for i in range(1, len(X) + 1):
        for j in range(1, len(X) + 1):
            if X[i - 1] == X[j - 1] and i != j:
                lookup[i][j] = lookup[i - 1][j - 1] + 1
            else:
                lookup[i][j] = max(lookup[i - 1][j], lookup[i][j - 1])
 
 
X = 'ATACTCGGA'

lookup = []
for y in range(len(X) + 1):
    t = []
    for x in range(len(X) + 1):
        t.append(0)
    lookup.append(t)


LRSLength(X, lookup)

print(LRS(X, len(X), len(X), lookup))

