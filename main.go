package main
import (
  "./tcp"
  "./serial"
  "time"
  "fmt"
  "strings"
  "log"
  "github.com/yosssi/gmq/mqtt"
  "github.com/yosssi/gmq/mqtt/client"
)

func main() {
  pretty_names := map[string]string{
    "10088D8D02080000": "Kühlschrank",
    "101C4D8D020800A5": "Tiefkühlfach",
    "10FF598D020800B6": "Küchenschrank",
  }

  server := tcp.New("[::]:1337")
  go server.Listen()

  mqtt_client := client.New(&client.Options{
    ErrorHandler: func(err error) {
      log.Println(err)
    },
  })
  defer mqtt_client.Terminate()
  err := mqtt_client.Connect(&client.ConnectOptions{
    Network:  "tcp",
    Address:  "172.25.7.119:1883",
    ClientID: []byte("kuehlschrank"),
  })
  if err != nil {
    log.Fatal(err)
  }

  ser := serial.Listen()
  for true {
    sensors := serial.QuerySensors(ser)
    newmessage := "Aktuelle Temperaturen in Lukas Villa: "
    var temps []string

    for _, s := range sensors {
      mqtt_client.Publish(&client.PublishOptions{
        QoS: mqtt.QoS0,
        TopicName: []byte(fmt.Sprintf("sensors/ds18b20/%s", s.Rom)),
        Message: []byte(fmt.Sprintf("%.3f", s.Temperature)),
      })
      if pretty_names[s.Rom] != "" {
        temps = append(temps, fmt.Sprintf("%s: %.2f°C", pretty_names[s.Rom], s.Temperature))
      }
    }
    newmessage += strings.Join(temps, ", ")
    newmessage += " ("
    newmessage += time.Now().Format(time.RFC3339)
    newmessage += ")"
    tcp.SetMessage(newmessage)

    log.Println("Sleeping 30 seconds...")
    time.Sleep(30 * time.Second)
  }
}
