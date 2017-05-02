<?php

    class Logger {

      private $socket;
      private $port;
      private $host;
      private $proto;
      private $format;
      private $debug = false;

      public function __construct($host = '127.0.0.1', $port = 5000, $proto = 'udp', $format='plain') {
        if (strtolower($proto) == 'udp') {
          if(!($this->socket = socket_create(AF_INET, SOCK_DGRAM, SOL_UDP))) {
              $errorcode = socket_last_error();
              $errormsg = socket_strerror($errorcode);
              die("Couldn't create UDP socket: [$errorcode] $errormsg \n");
          }
          $this->proto = 'udp';
        } else {
          if(!($this->socket = socket_create(AF_INET, SOCK_STREAM, 0))) {
              $errorcode = socket_last_error();
              $errormsg = socket_strerror($errorcode);
              die("Couldn't create TCP socket: [$errorcode] $errormsg \n");
          }
          $this->proto = 'tcp';
          $res = socket_connect($this->socket, $host, $port);
        }
        $this->host = $host;
        $this->port = $port;
        $this->format = ( in_array($format, ['plain','json']) ? $format : 'plain' );
      }

      public function setDebug() {
        $this->debug = true;
      }

      public function sendMsg($msg) {
        $payload = ($this->format == 'json' ? json_encode($msg) : $msg);
        if ($this->debug) {
          echo sprintf("Logging msg via %s %s:%d: %s\n", strtoupper($this->proto), $this->host, $this->port, $payload);
        }
        if ($this->proto == 'udp') {
          socket_sendto($this->socket, $payload, strlen($payload), 0, $this->host, $this->port);
        } else {
          socket_write($this->socket, $payload, strlen($payload));
        }
      }

      public function __destruct() {
        if ($this->debug) {
          echo "Closing socket...\n";
        }
        socket_close($this->socket);
      }
    }


    # Start a logger to accept default plain text messages over UDP
    $logger = new Logger('127.0.0.1', 6000);
    //$logger->setDebug(); // Uncomment for debugging
    $logger->sendMsg('This is a test plain text message from my test application');

    # Start a logger to accept default plain text messages over TCP
    /*
    $logger = new Logger('127.0.0.1', 6000, 'tcp', 'json');
    //$logger->setDebug(); // Uncomment for debugging
    $logger->sendMsg( 
      array(
        'type' => 'php_app_json',
        'message' => 'This is a test message from the PHP logger',
        'service' => 'payments',
        'log_level' => 'INFO'
      )
    );
    */
?>