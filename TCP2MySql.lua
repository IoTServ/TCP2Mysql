id = 'yourid'  --Device ID
ssid = 'yourSSID'  --WiFi Name
ssidpwd = 'WiFiPassword'  --view on http://www.mcunode.com/proxy/<yourid>/<anystring> like:http://www.mcunode.com/proxy/4567/index.html
ip="yourServerIpOrHostName"  --Eg："192.168.1.105" or“www.mcunode.com”
--设置前4行的参数：你自定义的id，wifi用户名，密码，你的服务器ip：服务器自己运行自己的版本，服务器上同时还得创建一个数据库名：mcunode，用户名：root 密码：root
function startServer()
print(wifi.sta.getip())
sk=net.createConnection(net.TCP, 0)
sk:on("receive", function(sck, c)
	node.input(c)
	end )   --print(c)
sk:on("connection", function(sck, c) 
--print(c)
sk:send(id)
tmr.alarm(2, 30000, 1, function() 
	sk:send('<h1></h1>')
end)
tmr.alarm(3, 3000, 1, function() 
	sk:send('wodechucun测试')
end)
end )
sk:on("disconnection",function(conn,c) 
         --node.output(nil) 
		 print('reconnect')
		 sk:connect(8002,ip)
		 sk:send(id)
      end)
sk:connect(8002,ip)
end
wifi.setmode(wifi.STATION)
wifi.sta.config(ssid,ssidpwd)    --set your ap info !!!!!!
wifi.sta.autoconnect(1)
tmr.alarm(1, 1000, 1, function() 
   if wifi.sta.getip()==nil then
      print("Connect AP, Waiting...") 
   else
      startServer()
      tmr.stop(1)
   end
end)