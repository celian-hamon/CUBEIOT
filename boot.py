# boot.py -- run on boot-up
import si7021
import machine
import time
import urequests as requests
from lcdi2c import LCDI2C
import json
import network

i2c = machine.I2C(scl=machine.Pin(2), sda=machine.Pin(0), freq=400000)
temp_sensor = si7021.Si7021(i2c)
lcd = LCDI2C(i2c, cols=2, rows=16)
lcd.backlight()

with open("config.json") as f:
    config = json.load(f)

station = network.WLAN(network.STA_IF)
station.active(True)
if (station.isconnected() != True):
    station.connect(config['ssid'], config['mdp'])
    if(station.isconnected() == True):
        print('Connected')
    else:
        print('Not Connected')
else :
    print('network config:', station.ifconfig())        

rtc = machine.RTC()

while True:
    for i in range( 20 ):
        Time = rtc.datetime()
        lcd.set_cursor((0,0))
        lcd.print('IP: 192.168.235')
        lcd.set_cursor((0,1))
        date = str(Time[2])+'/'+str(Time[1])+'/'+str(Time[0])+' - '+str(Time[4])+':'+str(Time[5])+':'+str(Time[6])
        lcd.print(date)
        lcd.set_cursor((24,0))
        lcd.print('Temp: ' + str(round(temp_sensor.temperature,2)))
        lcd.set_cursor((24,1))
        lcd.print('Hum : ' + str(round(temp_sensor.relative_humidity,2)))
        lcd.scroll_display()
        time.sleep( 1 )
        data = {'temperature': temp_sensor.temperature, 'humidite': temp_sensor.relative_humidity}
        r = requests.post('http://192.168.107.235:5000/releve', headers={'sonde': config['id']}, json = data)
        json_body = json.loads(r.text)
        print(json_body)
