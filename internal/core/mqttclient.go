package core

import (
    "fmt"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient handles MQTT messaging
type MQTTClient struct {
    client       mqtt.Client
    isConnected  bool
    clientID     string
    broker       string
    statusTopic  string
    connectToken mqtt.Token
}

// NewMQTTClient creates a new MQTT client with the specified configuration
func NewMQTTClient(clientID, broker, statusTopic string) *MQTTClient {
    return &MQTTClient{
        clientID:    clientID,
        broker:      broker,
        statusTopic: statusTopic,
    }
}

// Connect establishes a connection to the MQTT broker
func (m *MQTTClient) Connect() error {
    // Configure MQTT options
    opts := mqtt.NewClientOptions().
        AddBroker(m.broker).
        SetClientID(m.clientID).
        SetAutoReconnect(true).
        SetCleanSession(true).
        SetKeepAlive(60 * time.Second).
        SetPingTimeout(10 * time.Second)

    // Set connection callbacks
    opts.SetOnConnectHandler(m.onConnect)
    opts.SetConnectionLostHandler(m.onConnectionLost)

    // Create client
    m.client = mqtt.NewClient(opts)

    // Attempt connection
    m.connectToken = m.client.Connect()
    if m.connectToken.Wait() && m.connectToken.Error() != nil {
        return fmt.Errorf("failed to connect to MQTT broker: %v", m.connectToken.Error())
    }

    return nil
}

// Disconnect closes the connection to the MQTT broker
func (m *MQTTClient) Disconnect() {
    if m.client != nil && m.client.IsConnected() {
        m.client.Disconnect(1000) // 1 second timeout
    }
}

// IsConnected returns the current connection status
func (m *MQTTClient) IsConnected() bool {
    return m.isConnected
}

// SendMessage publishes a message to the configured status topic
func (m *MQTTClient) SendMessage(message string) error {
    if !m.IsConnected() {
        return fmt.Errorf("MQTT client not connected")
    }

    token := m.client.Publish(m.statusTopic, 0, false, message)
    if token.Wait() && token.Error() != nil {
        return fmt.Errorf("failed to publish MQTT message: %v", token.Error())
    }

    return nil
}

// SendStatusUpdate sends a service status message with timestamp
func (m *MQTTClient) SendStatusUpdate() error {
    statusMsg := fmt.Sprintf("Service status: running, timestamp: %s",
        time.Now().Format(time.RFC3339))
    return m.SendMessage(statusMsg)
}

// Connection event handlers
func (m *MQTTClient) onConnect(client mqtt.Client) {
    m.isConnected = true
}

func (m *MQTTClient) onConnectionLost(client mqtt.Client, err error) {
    m.isConnected = false
}