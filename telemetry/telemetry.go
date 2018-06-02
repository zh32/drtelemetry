package telemetry

import (
	"encoding/binary"
	"bytes"
	"net"
	"fmt"
	"flag"
)
/**
 * Protocol specification: https://docs.google.com/spreadsheets/d/1UTgeE7vbnGIzDz-URRk2eBIPc_LR1vWcZklp7xD9N0Y/edit#gid=0
 */
type TelemetryData struct {
	Time                   float32 `json:"time"`
	LapTime                float32 `json:"lapTime"`
	LapDistance            float32 `json:"lapDistance"`
	TotalDistance          float32 `json:"distance"`
	X                      float32 `json:"-"` // World space position
	Y                      float32 `json:"-"` // World space position
	Z                      float32 `json:"-"` // World space position
	Speed                  float32 `json:"speed"`
	Xv                     float32 `json:"-"` // Velocity in world space
	Yv                     float32 `json:"-"` // Velocity in world space
	Zv                     float32 `json:"-"` // Velocity in world space
	Xr                     float32 `json:"-"` // World space right direction
	Yr                     float32 `json:"-"` // World space right direction
	Zr                     float32 `json:"-"` // World space right direction
	Xd                     float32 `json:"-"` // World space forward direction
	Yd                     float32 `json:"-"` // World space forward direction
	Zd                     float32 `json:"-"` // World space forward direction
	Susp_pos_bl            float32 `json:"-"`
	Susp_pos_br            float32 `json:"-"`
	Susp_pos_fl            float32 `json:"-"`
	Susp_pos_fr            float32 `json:"-"`
	Susp_vel_bl            float32 `json:"-"`
	Susp_vel_br            float32 `json:"-"`
	Susp_vel_fl            float32 `json:"-"`
	Susp_vel_fr            float32 `json:"-"`
	Wheel_speed_bl         float32 `json:"-"`
	Wheel_speed_br         float32 `json:"-"`
	Wheel_speed_fl         float32 `json:"-"`
	Wheel_speed_fr         float32 `json:"-"`
	Throttle               float32 `json:"throttlePosition"`
	Steer                  float32 `json:"steerPosition"`
	Brake                  float32 `json:"brakePosition"`
	Clutch                 float32 `json:"clutchPosition"`
	Gear                   float32 `json:"gear"`
	Gforce_lat             float32 `json:"-"`
	Gforce_lon             float32 `json:"-"`
	Lap                    float32 `json:"-"`
	EngineRate             float32 `json:"rpm"`
	Sli_pro_native_support float32 `json:"-"` // SLI Pro support
	Car_position           float32 `json:"-"` // car race position
	Kers_level             float32 `json:"-"` // kers energy left
	Kers_max_level         float32 `json:"-"` // kers maximum energy
	Drs                    float32 `json:"-"` // 0 = off, 1 = on
	Traction_control       float32 `json:"-"` // 0 (off) - 2 (high)
	Anti_lock_brakes       float32 `json:"-"` // 0 (off) - 1 (on)
	Fuel_in_tank           float32 `json:"-"` // current fuel mass
	Fuel_capacity          float32 `json:"-"` // fuel capacity
	In_pits                float32 `json:"-"` // 0 = none, 1 = pitting, 2 = in pit area
	Sector                 float32 `json:"-"` // 0 = sector1, 1 = sector2 float32 `json:"-"` 2 = sector3
	Sector1_time           float32 `json:"-"` // time of sector1 (or 0)
	Sector2_time           float32 `json:"-"` // time of sector2 (or 0)
	Brakes_temp            float32 `json:"-"` // brakes temperature (centigrade)
	Wheels_pressure        float32 `json:"-"` // wheels pressure PSI
	Teainfo                float32 `json:"-"` // team ID
	Total_laps             float32 `json:"-"` // total number of laps in this race
	Track_size             float32 `json:"-"` // track size meters
	Last_lap_time          float32 `json:"-"` // last lap time
	Max_rpm                float32 `json:"-"` // cars max RPM, at which point the rev limiter will kick in
	Idle_rpm               float32 `json:"-"` // cars idle RPM
	Max_gears              float32 `json:"-"` // maximum number of gears
	SessionType            float32 `json:"-"` // 0 = unknown, 1 = practice, 2 = qualifying, 3 = race
	DrsAllowed             float32 `json:"trackLength"` // 0 = not allowed, 1 = allowed, -1 = invalid / unknown
	Track_number           float32 `json:"-"` // -1 for unknown, 0-21 for tracks
	VehicleFIAFlags        float32 `json:"maxRpm"` // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
}

var Addr = flag.String("udp", "localhost:10001", "udp service address")
var dataChannel = make(chan TelemetryData)

func ReadFromBytes(data []byte, telemetryData *TelemetryData) error {
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.LittleEndian, telemetryData); err != nil {
		return err
	}
	return nil
}

func RunServer(host string) (chan TelemetryData, chan struct{}) {
	quit := make(chan struct{})

	go func() {
		protocol := "udp"

		udpAddr, err := net.ResolveUDPAddr(protocol, host)
		if err != nil {
			fmt.Println("Wrong Address")
			return
		}

		udpConn, err := net.ListenUDP(protocol, udpAddr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Listening for telemetry data on %s\n", *Addr)
		defer udpConn.Close()

		for {
			select {
			default:
				HandleConnection(udpConn)
			case <-quit:
				fmt.Println("Stopping telemetry server")
				return
			}
		}
	}()

	return dataChannel, quit
}

func HandleConnection(conn *net.UDPConn) {
	buf := make([]byte, 264)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error Reading")
		fmt.Println(err)
		return
	} else {
		var reader = bytes.NewBuffer(buf[:n])
		var foo TelemetryData
		err := binary.Read(reader, binary.LittleEndian, &foo)
		if err != nil {
			fmt.Println(err)
		}
		dataChannel <- foo
	}
}
