# USR-PT170 mqtt信息

MQTT用户名：0893d019-c63b-48b3-b40
MQTT密码：c49bd8e
设备遥测上报主题 (更多上行主题请参考MQTT文档)
devices/telemetry

MQTT ClientID(自定义标识，需确保唯一性)
mqtt_112f3ddf-587

接入地址
47.92.253.145:1883


{ 
"Alarm": "1",                                      //"Alarm": "1"代表报警。无则为正常，温湿光的越限值任一满足，则触发报警。 
"sys_local_time": "2023-05-27,22:35:44",      //时间戳，表示本地时间 
"SN": "00501521042000019454",                  //设备 SN 
"temperature": "25",                      //温度值 
"humidity": "55",                            //湿度值 
"illuminance": "10000"}                //光照度