package main

import (
  "fmt"
  "os"
  "os/signal"
  "io/ioutil"
  "github.com/stratoberry/go-gpsd"
  "github.com/ptrv/go-gpx"
)

func main() {
  var gps *gpsd.Session
  var err error


  if gps, err = gpsd.Dial("localhost:2947"); err != nil {
    panic(fmt.Sprintf("Failed to connect to GPSD: ", err))
  }

  g := gpx.NewGpx()
  gpxTrack := gpx.Trk{}
  gpxSegment := gpx.Trkseg{}

  tpvfilter := func(r interface{}) {
    tpv := r.(*gpsd.TPVReport)
    gpxSegment.Points = append(gpxSegment.Points, gpx.Wpt{Timestamp: tpv.Time, Lat: tpv.Lat, Lon: tpv.Lon, Ele: tpv.Alt})
    //fmt.Println("TPV", tpv.Tag, tpv.Mode, tpv.Time, tpv.Lat, tpv.Lon, tpv.Alt, tpv.Epx, tpv.Epy, tpv.Epv)
  }

/*  skyfilter := func(r interface{}) {
    sky := r.(*gpsd.SKYReport)
    fmt.Println("SKY", len(sky.Satellites), "satellites")
  }
*/

  gps.AddFilter("TPV", tpvfilter)
//  gps.AddFilter("SKY", skyfilter)

c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt)
go func(){
  for sig := range c {
    gpxTrack.Segments = append(gpxTrack.Segments, gpxSegment)
    g.Tracks = append(g.Tracks, gpxTrack)
    err := ioutil.WriteFile("output.gpx", g.ToXML(), 0644)
    if err != nil {
      panic(err)
    }
    fmt.Println(sig)
    os.Exit(1)
  }
}()

  done := gps.Watch()
  <-done



}
