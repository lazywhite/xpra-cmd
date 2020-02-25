# -*- coding: utf-8 -*-
#
# Copyright © 2020 white <white@Whites-Mac-Air.local>
#
# Distributed under terms of the MIT license.

"""
"""

from flask import Flask, request, jsonify


app = Flask('xpra-cmd')


@app.route('/avlcloud/api/apps/', methods=['GET'])
def launch():
    cmds = [{"id":"8bea6556-8c98-4a7e-bfb6-bfafb0fef6a2","timestamp_created":"2019-12-05T06:00:19.490584Z","timestamp_updated":"2019-12-05T06:00:19.490690Z","label":"ffffffffffff","command":"fffffffffffff","category":"fffffffffffff","icon":"fffffffffffff","rank":0,"enabled":True},{"id":"df7fee2e-21fa-44c7-a766-cd5561e9bef9","timestamp_created":"2019-08-08T06:48:49.323178Z","timestamp_updated":"2019-08-08T06:48:49.323305Z","label":"aa","command":"aa","category":"aa","icon":"terminal.svg","rank":0,"enabled":True},{"id":"6f7bb7e7-8b66-4e93-b0e0-89320a66fbf5","timestamp_created":"2019-08-06T05:59:38.324101Z","timestamp_updated":"2019-08-06T06:21:35.513630Z","label":"aa","command":"aa","category":"aa","icon":"aa","rank":0,"enabled":True},{"id":"03fe4d10-ad52-4f73-b341-0c3aa927621d","timestamp_created":"2019-08-06T05:32:15.415005Z","timestamp_updated":"2019-08-06T05:57:18.524887Z","label":"string","command":"string","category":"string","icon":"string","rank":0,"enabled":True},{"id":"92ca0f0e-20ed-4d0b-a57a-50ba7f0136a5","timestamp_created":"2019-08-05T08:45:23.071190Z","timestamp_updated":"2019-08-05T08:45:23.071321Z","label":"ls","command":"ls","category":"ls","icon":"ls","rank":2,"enabled":True}]
    return jsonify(cmds)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)

