# return true value: "1 2 3"
# exp
#
# current_positions -> firstly what we predicted
# previous_poisitions -> return value (ground truth) -> current_positions
#
# movement=['L','R','S']
# ground truth="1 2 3"
# predicted="1 2 3"

import packet_pb2

def move_prediction(initial_position, movement):

    m1, m2, m3 = movement

    if initial_position == '1 2 3':
        if m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '1 2 3'
        elif m1 == 'R' and m2 == 'L' and m3 == 'S':
            return '2 1 3'
        elif m1 == 'R' and m2 == 'L' and m3 == 'L':
            return '2 3 1'
        elif m1 == 'R' and m2 == 'S' and m3 == 'L':
            return '3 2 1'
        elif m1 == 'R' and m2 == 'R' and m3 == 'L':
            return '3 1 2'
        elif m1 == 'S' and m2 == 'R' and m3 == 'L':
            return '1 3 2'

    elif initial_position == '1 3 2':
        if m1 == 'S' and m2 == 'R' and m3 == 'L':
            return '1 2 3'
        elif m1 == 'R' and m2 == 'L' and m3 == 'R':
            return '2 1 3'
        elif m1 == 'R' and m2 == 'L' and m3 == 'S':
            return '2 3 1'
        elif m1 == 'R' and m2 == 'L' and m3 == 'L':
            return '3 2 1'
        elif m1 == 'R' and m2 == 'S' and m3 == 'L':
            return '3 1 2'
        elif m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '1 3 2'

    elif initial_position == '2 1 3':
        if m1 == 'L' and m2 == 'R' and m3 == 'S':
            return '1 2 3'
        elif m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '2 1 3'
        elif m1 == 'R' and m2 == 'S' and m3 == 'L':
            return '2 3 1'
        elif m1 == 'R' and m2 == 'R' and m3 == 'L':
            return '3 2 1'
        elif m1 == 'S' and m2 == 'R' and m3 == 'L':
            return '3 1 2'
        elif m1 == 'L' and m2 == 'R' and m3 == 'L':
            return '1 3 2'

    elif initial_position == '2 3 1':
        if m1 == 'L' and m2 == 'R' and m3 == 'R':
            return '1 2 3'
        elif m1 == 'L' and m2 == 'S' and m3 == 'R':
            return '2 1 3'
        elif m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '2 3 1'
        elif m1 == 'S' and m2 == 'R' and m3 == 'L':
            return '3 2 1'
        elif m1 == 'L' and m2 == 'R' and m3 == 'L':
            return '3 1 2'
        elif m1 == 'L' and m2 == 'R' and m3 == 'S':
            return '1 3 2'

    elif initial_position == '3 2 1':
        if m1 == 'L' and m2 == 'S' and m3 == 'R':
            return '1 2 3'
        elif m1 == 'L' and m2 == 'L' and m3 == 'R':
            return '2 1 3'
        elif m1 == 'S' and m2 == 'L' and m3 == 'R':
            return '2 3 1'
        elif m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '3 2 1'
        elif m1 == 'L' and m2 == 'R' and m3 == 'S':
            return '3 1 2'
        elif m1 == 'L' and m2 == 'R' and m3 == 'R':
            return '1 3 2'

    elif initial_position == '3 1 2':
        if m1 == 'L' and m2 == 'L' and m3 == 'R':
            return '1 2 3'
        elif m1 == 'S' and m2 == 'L' and m3 == 'R':
            return '2 1 3'
        elif m1 == 'R' and m2 == 'L' and m3 == 'R':
            return '2 3 1'
        elif m1 == 'R' and m2 == 'L' and m3 == 'S':
            return '3 2 1'
        elif m1 == 'S' and m2 == 'S' and m3 == 'S':
            return '3 1 2'
        elif m1 == 'L' and m2 == 'S' and m3 == 'R':
            return '1 3 2'

    return 'invalid' #tell josh to correct this code, cannot put else

def generate_right_movement(initial_position, ground_truth):
    for m1 in "LRS":
        for m2 in "LRS":
            for m3 in "LRS":
                if ground_truth == move_prediction(initial_position, [m1,m2,m3]):
                    return [m1,m2,m3]
    return ["N", "N", "N"]

def generate_position_message(initial_position, predicted_movement, ground_truth):
    predicted_position = move_prediction(initial_position, predicted_movement)
    SUCCESS = "Your Position is Correct!"
    if predicted_position == ground_truth:
        return [SUCCESS,SUCCESS,SUCCESS]

    messages = ["","",""]
    for user in [0,1,2]:
        predicted_move = predicted_movement[user]
        right_move = generate_right_movement(initial_position, ground_truth)[user]
        appended_msg = f" (detected '{predicted_move}' instead of '{right_move}')"
        msg = ""

        if (predicted_move == right_move):
            messages[user] = SUCCESS
            continue
        elif (predicted_move == 'L' and right_move == 'R') or (predicted_move == 'R' and right_move == 'L'):
            msg = f"Oops! You're going the Wrong Way!"
        elif predicted_move == "S" and right_move == 'R':
            msg = f"Oops! You're not Moving Right Enough!"
        elif predicted_move == "S" and right_move == 'L':
            msg = f"Oops! You're not Moving Left Enough!"
        elif right_move == 'S':
            msg = f"Oops! You're not Keeping Still Enough!"

        msg += appended_msg
#         msg += f" -{predicted_position}"

        messages[user] = msg

    return messages

print(generate_position_message("1 2 3", ["S", "S", "R"], "1 3 2"))