
import socket 
import json

'''
Sample logging client that writes to the local instance of udplogbeat listening on the configured port
'''

class Logger:

    def __init__(self, host='127.0.0.1', port=5000):
        self.host = host
        self.port = port
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP
        self.debug = False

    def enable_debug(self):
        self.debug = True

    def send_message(self, fmt, es_type, msg):
        payload = "{}:{}:{}".format(fmt, es_type, msg)
        if self.debug:
            print "Sending payload: {}".format(payload)
        self.socket.sendto(payload.encode('utf-8'), (self.host, self.port))


l = Logger('127.0.0.1', 5000)
l.enable_debug()
l.send_message('plain', 'python_app', 'This is a sample plaintext message')
l.send_message('json', 'python_app_json', json.dumps({'message': 'This is a JSON encoded message', 'application': 'my_app', 'log_level': 'INFO'}))
