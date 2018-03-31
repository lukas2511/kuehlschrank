package main
import (
  "./tcp"
  "./serial"
  "time"
  "fmt"
  "strings"
  "log"
)

func main() {
  pretty_names := map[string]string{
    "10088D8D02080000": "Kühlschrank",
    "101C4D8D020800A5": "Tiefkühlfach",
    "10FF598D020800B6": "Küchenschrank",
  }

  server := tcp.New("[::]:1337")
  go server.Listen()

  ser := serial.Listen()
  for true {
    sensors := serial.QuerySensors(ser)
    newmessage := "Aktuelle Temperaturen in Lukas Villa: "
    var temps []string

    for _, s := range sensors {
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
