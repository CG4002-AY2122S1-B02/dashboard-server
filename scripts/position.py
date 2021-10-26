#change is represented as -1 for Left, 0 for Stay, 1 for right
#position is represented as ABC or 123

def move_predict(initial, change):
    initial = ["1","2","3"]
    change = [1, -1, -1]

    #Contain possible indexed positions after change i.e. if possibleA = [1, 2], the user at A could now be at index 1 or 2
    possibleA = []
    possibleB = []
    possibleC = []

    if change[0] == -1:
        print("error left guy can't move left")
    elif change[0] == 0:
        possibleA = [0]
    elif change[0] == 1:
        possibleA = [1, 2]

    if change[1] == -1:
        possibleB = [0]
    elif change[1] == 0:
        possibleB = [1]
    elif change[1] == 1:
        possibleB = [2]

    if change[2] == -1:
        possibleC = [0,1]
    elif change[2] == 0:
        possibleC = [2]
    elif change[2] == 1:
        print("error right guy can't move right")

    #check possible B first cos it narrows down the most

    #subtract possibleB from possibleA and possibleC
    if possibleB[0] in possibleA:
        possibleA.remove(possibleB[0])

    if possibleB[0] in possibleC:
        possibleC.remove(possibleB[0])

    #subtract possibleB and possibleA from each other
    if len(possibleA) > 1 and len(possibleC) > 1:
        if possibleA[0] in possibleC:
            possibleC.remove(possibleA[0])
        if possibleC[0] in possibleA:
            possibleA.remove(possibleC[0])

    final = ["X", "X", "X"]
    final[possibleA[0]] = initial[0]
    final[possibleB[0]] = initial[1]
    final[possibleC[0]] = initial[2]

    print(" ".join(final))

move_predict([],[])