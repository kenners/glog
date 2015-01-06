package main

import (
    "encoding/csv"
    "fmt"
    "github.com/docopt/docopt-go"
    "github.com/stratoberry/go-gpsd"
    "os"
    "strconv"
)

func f2s(input_num float64) string {
    return strconv.FormatFloat(input_num, 'f', -1, 64)
}

func main() {
    usage := `Glog.

Usage:
  glog [-h HOST] [-p PORT] [-o FILE]
  glog --help
  glog --version

Options:
  -h HOST --host=HOST    Host [default: localhost]
  -p PORT --port=PORT    Port [default: 2947]
  -o FILE --output=FILE    Output file [default: glog_output.csv]
  --help    Show this screen.
  --version    Show version.`

    args, _ := docopt.Parse(usage, nil, true, "Glog 0.1", false)

    var gps *gpsd.Session
    var err error

    host := fmt.Sprintf("%v:%v", args["--host"],  args["--port"])
    if gps, err = gpsd.Dial(host); err != nil {
        panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
    }

    csvfile, err := os.Create(fmt.Sprintf("%v", args["--output"]))
    if err != nil {
        panic(err)
    }
    defer csvfile.Close()
    writer := csv.NewWriter(csvfile)
    erra := writer.Write([]string{"timestamp", "lat", "lon", "alt", "speed", "track", "climb", "ept", "epy", "epx", "epv", "eps", "epd", "epc", "mode"})
    if erra != nil {
        panic(erra)
    }
    writer.Flush()

    tpvfilter := func(r interface{}) {
        tpv := r.(*gpsd.TPVReport)
        if tpv.Time != "" && tpv.Mode == 3 {
            err := writer.Write([]string{tpv.Time, f2s(tpv.Lat), f2s(tpv.Lon), f2s(tpv.Alt), f2s(tpv.Speed), f2s(tpv.Track), f2s(tpv.Climb), f2s(tpv.Ept), f2s(tpv.Epy), f2s(tpv.Epx), f2s(tpv.Epv), f2s(tpv.Eps), f2s(tpv.Epd), f2s(tpv.Epc), fmt.Sprintf("%v", tpv.Mode)})
            if err != nil {
                panic(err)
            }
            writer.Flush()
        }
    }
    gps.AddFilter("TPV", tpvfilter)
    done := gps.Watch()
    <-done
}
