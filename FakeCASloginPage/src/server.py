from flask import Flask, request
app = Flask(__name__)
@app.route('/', methods=['POST'])
def result():
    print(request.form['username'], request.form['password'])
    return 'Received !' # response to your request.