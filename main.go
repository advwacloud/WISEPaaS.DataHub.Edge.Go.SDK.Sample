package main

import (
	"fmt"
	"runtime"
	"time"

	SDK "github.com/advwacloud/WISEPaaS.DataHub.Edge.Go.SDK"
)

func main() {

	go func() {
		for {
			runtime.NumGoroutine()
			// fmt.Printf("goroutine num = %d\n", NumGoroutine)
			time.Sleep(5 * time.Second)
		}
	}()

	quit := make(chan bool)

	options := SDK.NewEdgeAgentOptions()
	options.NodeID = "7654a6d2-7d1a-4b56-b397-f49555bc4160"
	options.ConnectType = SDK.ConnectType["DCCS"]
	options.DCCS.Key = "3528b8a09d6314169e200e412b588d4r"
	options.DCCS.URL = "https://api-dccs-ensaas.sa.wise-paas.com/"
	options.DataRecover = true

	// options := SDK.NewEdgeAgentOptions()
	// options.ConnectType = SDK.ConnectType["MQTT"]
	// options.DataRecover = true
	// options.MQTT.HostName = "127.0.0.1"
	// options.MQTT.Port = 1883

	interval := 1
	var timer chan bool = nil

	agent := SDK.NewAgent(options)
	agent.SetOnConnectHandler(func(a SDK.Agent) {
		fmt.Println("connect successfully")

		config := generateConfig()
		action := SDK.Action["Create"]
		_ = agent.UploadConfig(action, config)

		status := generateDeviceStatus()
		_ = agent.SendDeviceStatus(status)

		timer = setInterval(func() {
			data := generateData()
			ok := agent.SendData(data)
			if ok {
				fmt.Println(data)
			}
		}, interval, true)
	})
	agent.SetOnDisconnectHandler(func(a SDK.Agent) {
		fmt.Println("disconnect successfully")
	})
	agent.SetOnMessageReceiveHandler(func(args SDK.MessageReceivedEventArgs) {
		msgType := args.Type
		message := args.Message
		switch msgType {
		case SDK.MessageType["WriteValue"]: // message format: WriteDataMessage
			for _, device := range message.(SDK.WriteDataMessage).DeviceList {
				fmt.Println("DeviceId: ", device.ID)
				for _, tag := range device.TagList {
					fmt.Println("TagName: ", tag.Name, ", Value: ", tag.Value)
				}
			}
		case SDK.MessageType["ConfigAck"]: // message format: ConfigAckMessage
			fmt.Println(message.(SDK.ConfigAckMessage).Result)
		case SDK.MessageType["TimeSync"]: //message format: TimeSyncMessage
			fmt.Println(message.(SDK.TimeSyncMessage).UTCTime)
		}
	})

	err := agent.Connect()
	if err != nil {
		fmt.Println(err)
	}
	<-quit
}

func generateConfig() SDK.EdgeConfig {
	nodeConfig := generateNodeConfig()
	edgeConfig := SDK.EdgeConfig{
		Node: nodeConfig,
	}
	return edgeConfig
}

func generateNodeConfig() SDK.NodeConfig {
	var deviceNum = 1

	nodeConfig := SDK.NewNodeConfig()
	nodeConfig.SetType(SDK.EdgeType["Gateway"])

	for idx := 0; idx < deviceNum; idx++ {
		config := generateDeviceConfig(idx + 1)
		nodeConfig.DeviceList = append(nodeConfig.DeviceList, config)
	}
	return nodeConfig
}

func generateDeviceConfig(idx int) SDK.DeviceConfig {
	var deviceID = fmt.Sprintf("%s%d", "Device", idx)
	var analogNum = 3
	var discreteNum = 2
	var textNum = 1

	deviceConfig := SDK.NewDeviceConfig(deviceID)
	deviceConfig.SetName(fmt.Sprintf("%s%d", "Device", idx))
	deviceConfig.SetType("Smart Device")
	deviceConfig.SetDescription(fmt.Sprintf("%s %d", "Device ", idx))

	for idx := 0; idx < analogNum; idx++ {
		config := generateAnalogConfig(idx + 1)
		deviceConfig.AnalogTagList = append(deviceConfig.AnalogTagList, config)
	}
	for idx := 0; idx < discreteNum; idx++ {
		config := generateDiscreteConfig(idx + 1)
		deviceConfig.DiscreteTagList = append(deviceConfig.DiscreteTagList, config)
	}
	for idx := 0; idx < textNum; idx++ {
		config := generateTextConfig(idx + 1)
		deviceConfig.TextTagList = append(deviceConfig.TextTagList, config)
	}
	return deviceConfig
}

func generateAnalogConfig(idx int) SDK.AnalogTagConfig {
	var tagName = fmt.Sprintf("%s%d", "ATag", idx)

	config := SDK.NewAnaglogTagConfig(tagName)
	config.SetDescription(fmt.Sprintf("%s %d", "ATag", idx))
	config.SetReadOnly(false)
	config.SetArraySize(0)
	config.SetSpanHigh(1000.0)
	config.SetSpanLow(0.0)
	config.SetEngineerUnit("")
	config.SetIntegerDisplayFormat(4)
	config.SetFractionDisplayFormat(2)

	return config
}

func generateDiscreteConfig(idx int) SDK.DiscreteTagConfig {
	var tagName = fmt.Sprintf("%s%d", "DTag", idx)

	config := SDK.NewDiscreteTagConfig(tagName)
	config.SetDescription(fmt.Sprintf("%s %d", "DTag ", idx))
	config.SetArraySize(0)
	config.SetReadOnly(true)
	config.SetState0("No")
	config.SetState1("Yes")

	return config
}

func generateTextConfig(idx int) SDK.TextTagConfig {
	var tagName = fmt.Sprintf("%s%d", "TTag", idx)

	config := SDK.NewTextTagConfig(tagName)
	config.SetDescription(fmt.Sprintf("%s %d", "TTag", idx))
	config.SetReadOnly(false)
	config.SetArraySize(0)

	return config
}

var numF float64
var numI int

func generateData() SDK.EdgeData {
	numF += 1
	deviceNum := 1
	msg := SDK.EdgeData{
		Timestamp: time.Now(),
	}

	for idx := 0; idx < deviceNum; idx++ {
		analogNum := 3
		discreteNum := 2
		textNum := 1
		deviceID := fmt.Sprintf("%s%d", "Device", idx+1)
		for num := 0; num < analogNum; num++ {
			t := SDK.EdgeTag{
				DeviceID: deviceID,
				TagName:  fmt.Sprintf("%s%d", "ATag", num+1),
				Value:    numF,
			}

			//fmt.Println(rand.Float64())

			msg.TagList = append(msg.TagList, t)
		}
		for num := 0; num < discreteNum; num++ {
			t := SDK.EdgeTag{
				DeviceID: deviceID,
				TagName:  fmt.Sprintf("%s%d", "DTag", num+1),
				Value:    numI,
			}
			numI += 1
			msg.TagList = append(msg.TagList, t)
		}
		for num := 0; num < textNum; num++ {
			t := SDK.EdgeTag{
				DeviceID: deviceID,
				TagName:  fmt.Sprintf("%s%d", "TTag", num+1),
				Value:    fmt.Sprintf("%s%d", "Test", numI),
			}
			msg.TagList = append(msg.TagList, t)
		}
	}
	return msg
}

func setInterval(someFunc func(), seconds int, async bool) chan bool {
	interval := time.Duration(seconds) * time.Second
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				close(clear)
				return
			}
		}
	}()
	return clear
}

func generateDeviceStatus() SDK.EdgeDeviceStatus {
	status := SDK.EdgeDeviceStatus{
		Timestamp: time.Now(),
	}
	deviceNum := 1

	for idx := 0; idx < deviceNum; idx++ {
		s := SDK.DeviceStatus{
			ID:     fmt.Sprintf("%s%d", "Device", idx+1),
			Status: SDK.Status["Online"],
		}
		status.DeviceList = append(status.DeviceList, s)
	}
	return status
}
