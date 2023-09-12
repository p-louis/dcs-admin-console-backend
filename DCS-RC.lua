local DCSRC = {}
local socket = require("socket")
local JSON = loadfile("Scripts\\JSON.lua")()

local tcpServer
local currentMission = ""
local data
local client
local err

local tcpPort = 50051

function DCSRC.onSimulationStart()
  currentMission = DCS.getMissionName()

  if tcpServer == nil then
    local successful
    tcpServer = socket.tcp()
    successful, err = tcpServer:bind("127.0.0.1", tcpPort)
    tcpServer:listen(1)
    tcpServer:settimeout(0)
    if not successful then
      log.write("DCS-RC", log.ERROR, "Error opening tcp socket - "..tostring(err))
    else
      log.write("DCS-RC", log.INFO, "Opened connection")
    end
  end
end

function DCSRC.onSimulationFrame()
  if client ~= nil then
    client:settimeout(10)
    data, err = client:receive()

    if data ~= nil then
      local keys = JSON:decode(data)
      log.write("DCS-RC", log.INFO, "Received "..data)

      if keys["command"] == "load_mission" then
        net.load_mission(keys["mission_name"])
      elseif keys["command"] == "append_mission" then
        net.missionlist_append(keys["mission_name"])
      elseif keys["command"] == "delete_mission" then
        net.missionlist_delete(keys["index"])
      elseif keys["command"] == "run_mission" then
        net.missionlist_run(keys["index"])
      elseif keys["command"] == "clear_missionlist" then
        net.missionlist_clear()
      elseif keys["command"] == "get_mission" then
        client:send("{\"filename\": \""..currentMission.."\"}\n")
      elseif keys["command"] == "get_missionlist" then
        client:send("{\"missionlist\":" .. JSON:encode(net.missionlist_get()) .."}\n")
      elseif keys["command"] == "pause_mission" then
        DCS.setPause(true)
      elseif keys["command"] == "unpause_mission" then
        DCS.setPause(false)
      elseif keys["command"] == "get_pause" then
        local pauseState = DCS.getPause()
        client:send("{\"pause_state\": "..tostring(pauseState).."}\n")
      end
    end
  end
  client, err = tcpServer:accept()
  if client ~= nil then
    log.write("DCS-RC", log.INFO, "Client Connected")
  end
end

DCS.setUserCallbacks(DCSRC)  -- here we set our callbacks

local webCallbacks = {}
webCallbacks.onWebServerRequest = onWebServerRequest
DCS.setUserCallbacks(webCallbacks)

