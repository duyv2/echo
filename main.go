/*
xmpp_echo is a demo client that connect on an XMPP server and echo message received back to original sender.
*/

package main

import (
        "fmt"
        "log"
        "os"

        "gosrc.io/xmpp"
        "gosrc.io/xmpp/stanza"
)

func main() {
        config := xmpp.Config{
                TransportConfiguration: xmpp.TransportConfiguration{
                        Address: "ejabber.wchat.vn:5222",
                },
                Jid:          "echo@ejabber.wchat.vn/mobile",
                Credential:   xmpp.Password("123456"),
                StreamLogger: os.Stdout,
                //Insecure:     false,
                //TLSConfig: tls.Config{InsecureSkipVerify: true},
        }

        router := xmpp.NewRouter()
        router.HandleFunc("message", handleMessage)

        client, err := xmpp.NewClient(&config, router, errorHandler)
        if err != nil {
                log.Fatalf("%+v", err)
        }

        // If you pass the client to a connection manager, it will handle the reconnect policy
        // for you automatically.
        cm := xmpp.NewStreamManager(client, nil)
        log.Fatal(cm.Run())
}

func handleMessage(s xmpp.Sender, p stanza.Packet) {
        msg, ok := p.(stanza.Message)
        if !ok {
                _, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
                return
        }

        _, _ = fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", msg.Body, msg.From)
        reply := stanza.Message{Attrs: stanza.Attrs{To: msg.From}, Body: msg.Body}
        _ = s.Send(reply)
}

func errorHandler(err error) {
        fmt.Println(err.Error())
}