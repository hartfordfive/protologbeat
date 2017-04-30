<?php

    class Logger {

      private $socket;
      private $port;
      private $host;
      private $debug = false;

      public function __construct($host = '127.0.0.1', $port = 5000) {
        $this->socket = socket_create(AF_INET, SOCK_DGRAM, SOL_UDP);
        $this->host = $host;
        $this->port = $port;
      }

      public function setDebug() {
        $this->debug = true;
      }

      public function sendMsg($format, $esType, $msg) {
        $payload = $format . ':' . $esType . ':' . $msg;
        if ($this->debug) {
          echo "Sending payload: \"$payload\"\n";
        }
        socket_sendto($this->socket, $payload, strlen($payload), 0, $this->host, $this->port);
      }

      public function __destruct() {
        if ($this->debug) {
          echo "Closing socket...\n";
        }
        socket_close($this->socket);
      }
    }


    $logger = new Logger();
    $logger->setDebug();
    $logger->sendMsg('plain', 'php_app', 'This is a test message from the PHP logger');
    $logger->sendMsg('json', 'php_app_json', 
      json_encode( array(
        'message' => 'This is a test message from the PHP logger',
        'application' => 'php_app_json',
        'log_level' => 'INFO'
      )
    ));
?>