
import socket 
import json

'''
Sample logging client that writes to the local instance of protologbeat listening on the configured host/port/protocol.
'''

class Logger:

    def __init__(self, host='127.0.0.1', port=5000, proto='udp', msg_format='plain'):
        self.host = host
        self.port = port 
        if proto == 'udp':
          self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        else:
          self.proto = 'tcp'
          self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM) # TCP
        self.msg_format = msg_format
        if self.msg_format not in ['plain','json']:
          self.msg_format = 'plain'
        self.debug = False

    def send_message(self, msg):
        if self.msg_format == 'json':
          payload = json.dumps(msg)
        else:
          payload = msg
        self.socket.sendto(payload.encode('utf-8'), (self.host, self.port))


# Initializing udp connection and sending a plaintext message
l = Logger('127.0.0.1', 6000)
l.send_message('This is a sample plaintext message to be sent via udp')

# Initializing tcp connection and sending a json-encoded message
l = Logger('127.0.0.1', 6000, 'tcp', 'json')
l.send_message({'message': 'This is a JSON encoded message', 'type': 'python_app_json', 'application': 'my_app', 'log_level': 'INFO'})
