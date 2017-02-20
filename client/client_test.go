package client

import (
	"fmt"
	"testing"
)

var serverHost string = "ws://localhost:9000/websocket"

func getClient() *Client {
	ws := NewWebSocketClient()
	client := New(ws)
	return client
}

func getAuthClient() *Client {
	client := getClient()
	err := client.Connect(serverHost)
	if err != nil {
		fmt.Println("Failed to connect client")
		return nil
	}

	userId, err := client.Register("user", "pass")
	if err != nil {
		fmt.Println("Failed to register client")
		return nil
	}

	err = client.Authenticate(userId, "pass")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return client
}

func TestConnect(t *testing.T) {
	client := getClient()
	defer client.Disconnect()
	err := client.Connect(serverHost)
	if err == nil {
		fmt.Println("connected")
	} else {
		fmt.Println(err)
		t.Fail()
	}
}

func TestConnectAsync(t *testing.T) {
	client := getClient()
	defer client.Disconnect()
	ch := client.ConnectAsync(serverHost)
	connected := <-ch
	if connected {
		fmt.Println("connected")
	} else {
		fmt.Println("Failed to connect")
		t.Fail()
	}
}

func TestDisconnect(t *testing.T) {
	client := getClient()
	client.Connect(serverHost)
	client.Disconnect()
}

func TestDisconnectWithoutConnect(t *testing.T) {
	client := getClient()
	client.Disconnect()
}

func TestReconnect(t *testing.T) {
	client := getClient()
	defer client.Disconnect()
	err := client.Connect(serverHost)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = client.Reconnect()
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println("reconnected")
}

func TestReconnectNoHost(t *testing.T) {
	client := getClient()
	defer client.Disconnect()

	err := client.Reconnect()
	if err != nil {
		t.Log("no host")
		return
	} else {
		t.Fatal("Can't call reconnect without any connect")
	}
}

func TestReconnectAsync(t *testing.T) {
	client := getClient()
	defer client.Disconnect()

	err := client.Connect(serverHost)
	if err != nil {
		t.Fatal(err)
		return
	}

	ch := client.ReconnectAsync()
	reconnected := <-ch
	if reconnected {
		t.Log("reconnected")
	} else {
		t.Fatal("Failed to reconnect async")
	}
}

func TestReconnectAsyncNoHost(t *testing.T) {
	client := getClient()
	defer client.Disconnect()

	ch := client.ReconnectAsync()
	reconnected := <-ch
	if reconnected {
		t.Fatal("Need error for no host")
	} else {
		t.Log("no host")
	}
}

func TestRegister(t *testing.T) {
	client := getClient()
	defer client.Disconnect()

	client.Connect(serverHost)
	_, err := client.Register("user", "pass")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("registered")
	}
}

func TestAuthenticate(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	if client == nil {
		t.Fatal("Failed to authenticate")
	} else {
		t.Log("authenticated")
	}
}

func TestListContactsEmpty(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	contacts, err := client.ListContacts()
	if err != nil {
		t.Fatal(err)
	}

	if len(contacts) == 0 {
		t.Log("empty contacts")
	}
}

func TestAddContact(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	err := client.AddContact("user2")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("added")
	}
}

func TestListContactsOneContact(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	err := client.AddContact("user2")
	if err != nil {
		t.Fatal(err)
	} else {
		contacts, err := client.ListContacts()
		if err != nil {
			t.Fatal(err)
		}

		contactsLen := len(contacts)
		if contactsLen == 0 {
			t.Fatal("Empty contacts")
		}

		t.Log(contacts[0])
	}
}

func TestSendMessage(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	client.SendMessage(client.UserId(), "hello")
	msg := <-client.Out
	data := msg.String()
	if data != "hello" {
		t.Fatal("Wrong message delivered")
	} else {
		t.Log(data)
	}
}

func TestSendTextFile(t *testing.T) {
	client := getAuthClient()
	defer client.Disconnect()

	uri, err := client.SendTextFile(client.UserId(), "./client.go")
	if err != nil {
		t.Fatal("Failed to send file", err)
	} else {
		t.Log(uri)
	}
}

func ExampleClient_StartInteractiveMode() {
	ws := NewWebSocketClient()
	c := New(ws)
	err := c.Connect("ws://127.0.0.1:9000/websocket")
	fmt.Println("connected", err)
	userId, err := c.Register("pesho", "123")
	fmt.Println("registered", err)

	err = c.Authenticate(userId, "123")
	fmt.Println("authenticated", err)

	if err == nil {
		c.StartInteractiveMode()
	}
}
