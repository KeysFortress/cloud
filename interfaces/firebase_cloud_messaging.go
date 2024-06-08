package interfaces

type FirebaseCloudMessaging interface {
	DeviceSet(token string, deviceId string)
	DeviceTokenUpdated(token string, deviceId string)
	Send()
}
