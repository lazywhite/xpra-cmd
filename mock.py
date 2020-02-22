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
    cmds = ["ls", "sleep 3"]
    return jsonify(cmds)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)

