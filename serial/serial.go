package serial
import (
  "bytes"
  "log"
  "encoding/json"
  "time"
  "github.com/tarm/serial"
)

type Sensor struct {
  Rom string
  Temperature float32
}

func Listen() *serial.Port {
  ser_conf := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Second * 2}
  ser, err := serial.OpenPort(ser_conf)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Giving Arduino 5 seconds to boot...")
  time.Sleep(5 * time.Second)
  return ser
}

func QuerySensors(ser *serial.Port) []Sensor {
  log.Println("Querying sensors...")
  ser.Write([]byte("g\n"))

  rawdata := make([]byte, 2048)
  totalbytes := 0
  for true {
    tmpbuf := make([]byte, 128)
    n, err := ser.Read(tmpbuf)
    if err != nil || n == 0 {
      break
    }
    rawdata = bytes.Join([][]byte{rawdata[:totalbytes], tmpbuf[:n]}, []byte(""))
    totalbytes += n
  }
  rawdata = bytes.TrimSpace(rawdata)
  rawdata = bytes.Replace(rawdata, []byte("{"), []byte("[{"), 1)
  rawdata = bytes.Replace(rawdata, []byte("}"), []byte("},"), -1)
  rawdata[len(rawdata)-1] = 0x5D

  var sensors []Sensor
  err := json.Unmarshal(rawdata, &sensors)
  if err != nil {
    log.Println("Can't decode JSON data:", err)
  }
  log.Printf("%+v", sensors)

  time.Sleep(1 * time.Second)
  return sensors
}
