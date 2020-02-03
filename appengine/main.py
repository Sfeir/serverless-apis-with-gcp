# [START gae_python37_app]
from flask import Flask
from flask import request

import logging

from google.cloud import datastore

client = datastore.Client()

ANCESTOR_KIND = 'Environment';
KIND = 'Microservice';

# If `entrypoint` is not defined in app.yaml, App Engine will look for an app
# called `app` in `main.py`.
app = Flask(__name__)


@app.route('/', methods=["GET"])
def appInventory():
    year = request.args.get('year')

    ancestor = client.key(ANCESTOR_KIND, 'App Engine')
    query = client.query(kind=KIND, ancestor=ancestor)

    if year != None :
      try:
        query.add_filter('year', '=', int(year)) 
      except ValueError as ex:
        print('"%s" cannot be converted to an int' % year)

    result = '<h1>App Engine Applications Inventory</h1><ul>';
    query_iter = query.fetch()
    for microservice in query_iter:
      result += '<li>' + microservice['name'] + ' (' + str(microservice['year']) + ')</li>'
    result += '</ul>'

    return result


if __name__ == '__main__':
    # This is used when running locally only. When deploying to Google App
    # Engine, a webserver process such as Gunicorn will serve the app. This
    # can be configured by adding an `entrypoint` to app.yaml.
    app.run(host='127.0.0.1', port=8080, debug=True)
# [END gae_python37_app]
