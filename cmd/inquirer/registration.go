/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2012  Radu Ioan Fericean

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package main

import (
	"os/signal"
	"fmt"
	"log"
	"net/rpc"
	"net/http"
	"os"
	"syscall"
	"time"
)

/*
RPC Server that handles the registering and unregistering of raters.
*/
type RaterServer byte

func listenToRPCRaterRequests(){
	raterServer := new(RaterServer)
	rpc.Register(raterServer)
	rpc.HandleHTTP()
	http.ListenAndServe(*raterAddress, nil)
}

/*
Listens for SIGTERM, SIGINT, SIGQUIT system signals and shuts down all the registered raters.
*/
func StopSingnalHandler() {
	log.Print("Handling stop signals...")
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sig := <-c
	log.Printf("Caught signal %v, sending shutdownto raters\n", sig)
	var reply string
	for i, client := range raterList.clientConnections {
		client.Call("Storage.Shutdown", "", &reply)
		log.Printf("Shutdown rater %v: %v ", raterList.clientAddresses[i], reply)
	}
	os.Exit(1)
}

/*
RPC method that receives a rater address, connects to it and ads the pair to the rater list for balancing
*/
func (rs *RaterServer) RegisterRater(clientAddress string, replay *byte) error {
	time.Sleep(1 * time.Second) // wait a second for Rater to start serving
	client, err := rpc.Dial("tcp", clientAddress)
	if err != nil {
		log.Print("Could not connect to client!")
		return err
	}
	raterList.AddClient(clientAddress, client)
	log.Print(fmt.Sprintf("Rater %v registered succesfully.", clientAddress))
	return nil
}

/*
RPC method that recives a rater addres gets the connections and closes it and removes the pair from rater list.
*/
func (rs *RaterServer) UnRegisterRater(clientAddress string, replay *byte) error {
	client, ok := raterList.GetClient(clientAddress)
	if ok {
		client.Close()
		raterList.RemoveClient(clientAddress)
		log.Print(fmt.Sprintf("Rater %v unregistered succesfully.", clientAddress))
	} else {
		log.Print(fmt.Sprintf("Server %v was not on my watch!", clientAddress))
	}
	return nil
}
