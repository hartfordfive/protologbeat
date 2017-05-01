Logger = {}
Logger.__index = Logger

function Logger.init(host,port,proto,format)
  local sock = require("socket") 
  local json = require("cjson")
  local lgr = {}             -- our new object
  setmetatable(lgr, Logger)  -- make Account handle lookup

  if proto == "tcp" then
    lgr.socket = sock.tcp()      -- initialize our object
  else
    lgr.socket = sock.udp()
  end
  lgr.socket:settimeout(0)
  lgr.host = host
  lgr.port = port
  if lgr.format == 'json' then
    lgr.format = 'json'
  else
    lgr.format = 'plain'
  end
  return lgr
end

function Logger:sendMsg(msg)
  local payload
  if self.format == 'json' then
    payload = self.json.encode(msg)
  else
    payload = msg
  end
  self.socket:sendto(payload, self.host, self.port)
end

-- Start logger client to send plain-text formated message to protologbeat listening on UDP host/port
logger = Logger.init('127.0.0.1', 6000, "udp", "plain")
logger:sendMsg('This is a sample message sent from the Lua logger.')

-- Start logger client to send json formated message to protologbeat listening on TCP host/port
--logger = Logger.init('127.0.0.1', 6000, "tcp", "json")
--logger:sendMsg({type = 'lua_app_json', message = 'This is a sample message sent from the Lua logger.', log_level = 'INFO'})
