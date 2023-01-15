from flask import Flask, request, render_template

app = Flask(__name__)
@app.route('/', methods=['POST'])
def result():
  print(request.form['username'], request.form['password'])
  return render_template('redirect.html')
