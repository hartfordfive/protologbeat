Logger = {}
Logger.__index = Logger

function Logger.init(host,port)
  local sock = require("socket") 
  local lgr = {}             -- our new object
  setmetatable(lgr, Logger)  -- make Account handle lookup
  lgr.socket = sock.udp()      -- initialize our object
  lgr.socket:settimeout(0)
  lgr.host = host
  lgr.port = port
  return lgr
end

function Logger:sendMsg(format, esType, msg)
  payload = format .. ":" .. esType .. ":" .. msg 
  self.socket:sendto(payload, self.host, self.port)
end

logger = Logger.init('127.0.0.1', 5000)
logger:sendMsg('plain','lua_app','This is a sample message sent from the Lua logger.')
