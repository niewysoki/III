from flask import Flask, request, render_template

users = {}

app = Flask(__name__)

@app.route('/', methods=['POST'])
def result():
  print(request.form['username'], request.form['password'])
  return render_template('redirect.html')

@app.route('/', methods=['GET'])
def result_get():
  return render_template('page.html')