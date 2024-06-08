package implementations

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"leanmeal/api/repositories"
)

type FirebaseCloudMessaging struct {
	AccessKeysRepository repositories.AccessKeysRepository
	ServerKey            string
}

func (f *FirebaseCloudMessaging) DeviceSet(token string, deviceId string) {

}

func (f *FirebaseCloudMessaging) DeviceTokenUpdated(token string, deviceId string) {

}

func (f *FirebaseCloudMessaging) Send() {
	// Path to your Firebase service account key file
	sa := option.WithCredentialsFile(f.ServerKey)

	// Create a Firebase app instance
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Get a messaging client from the app
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// Define the message to send
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Hello",
			Body:  "World",
		},
		Token: "esQ5A7nkQeaC_g_qy1HRWD:APA91bF9gNVEJr9i-NCS65gFvJjNCKioHp-M86QCoEnuj6p75dwBnoR45kQtWg-LOpg85yqT4IqeWiGolR3x5-gFxl2xKw65r1wrje8QEtPEBXBAit6kSjsUZHISZWh8y_dKzx6gKCZQ",
	}

	// Send the message via FCM
	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
	}

	// Print the response from FCM
	fmt.Printf("Successfully sent message: %s\n", response)
}
