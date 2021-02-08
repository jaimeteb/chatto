from flask import Flask, Response, request, jsonify

app = Flask(__name__)

def wrong_option(data):
    return {
        "fsm": {
            "state": data.get("db").get("question_1"),
            "slots": data.get("fsm").get("slots"),
        },
        "res": "Select one of the options."
    }

def validate_ans_1(data: dict) -> dict:
    if data.get("txt") not in ["1", "2", "3"]:
        return jsonify(wrong_option(data))

    return jsonify({
        "fsm": data.get("fsm"),
        "res": "Question 2:\n" +
            "What is the capital of the state of Utah?\n" +
            "1. Salt Lake City\n" +
            "2. Jefferson City\n" +
            "3. Cheyenne",
    })

def validate_ans_2(data: dict) -> dict:
    if data.get("txt") not in ["1", "2", "3"]:
        return jsonify(wrong_option(data))

    return jsonify({
        "fsm": data.get("fsm"),
        "res": "Question 3:\n" +
			"Who painted Starry Night?\n" +
			"1. Pablo Picasso\n" +
			"2. Claude Monet\n" +
			"3. Vincent Van Gogh",
    })

def score(data: dict) -> dict:
    if data.get("txt") not in ["1", "2", "3"]:
        return jsonify(wrong_option(data))
    
    slots = data.get("fsm").get("slots", {})
    answer_1 = slots.get("answer_1")
    answer_2 = slots.get("answer_2")
    answer_3 = slots.get("answer_3")
    
    score = 0
    score = score + 1 if answer_1 == "2" else score
    score = score + 1 if answer_2 == "1" else score
    score = score + 1 if answer_3 == "3" else score

    message = ""
    if score == 0:
        message = "You got 0/3 answers right.\nBetter luck next time!"
    elif score == 1:
        message = "You got 1/3 answers right.\nKeep trying!"
    elif score == 2:
        message = "You got 2/3 answers right.\nPretty good!"
    elif score == 3:
        message = "You got 3/3 answers right.\nYou are good! Congrats!"

    return jsonify({
        "fsm": data.get("fsm"),
        "res": message,
    })

func_map = {
    "ext_val_ans_1": validate_ans_1,
    "ext_val_ans_2": validate_ans_2,
    "ext_score": score,
}


@app.route("/ext/get_all_funcs", methods=["GET"])
def get_all_funcs():
    return jsonify(list(func_map.keys()))

@app.route("/ext/get_func", methods=["POST"])
def get_func():
    data = request.get_json()
    req = data.get("req")
    f = func_map.get(req)
    if not f:
        return Response(status=400)
    return f(data)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8770)
