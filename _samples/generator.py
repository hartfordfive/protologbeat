from logger import Logger
import time, random

l = Logger('127.0.0.1', 6000, 'udp', 'json')

for i in range(10000):
  l.send_message({'message': 'This is JSON encoded message #{}'.format(i), 'type': 'generator_test', 'id': int(i), 'log_level': 'INFO'})
  time.sleep(random.uniform(0.0001, 0.0010))