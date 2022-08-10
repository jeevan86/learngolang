package device

import (
	"fmt"
	"github.com/google/gopacket/pcap"
)

func printDevice(devices []pcap.Interface) {
	logger.Info("Devices found:")
	for _, device := range devices {
		printDeviceInfo(device)
	}
}

func printDeviceInfo(device pcap.Interface) {
	logger.Info(fmt.Sprintf("\nName: %s", device.Name))
	logger.Info(fmt.Sprintf("Description: %s", device.Description))
	logger.Info(fmt.Sprintf("Devices addresses: %s", device.Description))
	for _, address := range device.Addresses {
		printAddressInfo(address)
	}
}

func printAddressInfo(address pcap.InterfaceAddress) {
	logger.Info(fmt.Sprintf("- IP address: %s", address.IP))
	logger.Info(fmt.Sprintf("- Subnet mask: %s", address.Netmask))
	logger.Info(fmt.Sprintf("- Broad address: %s", address.Broadaddr))
	logger.Info(fmt.Sprintf("- P2P: %s", address.P2P))
}
