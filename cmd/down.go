package cmd

import (
	"fmt"
	"log"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop your application",
	Long:  ``,
	Run:   down,
}

func init() {
	RootCmd.AddCommand(downCmd)
}

func down(cmd *cobra.Command, args []string) {

	//read the harbor compose file
	var harborCompose = DeserializeHarborCompose(HarborComposeFile)

	//read the docker compose file
	var dockerCompose = DeserializeDockerCompose(DockerComposeFile)
	if dockerCompose.Version != "2" {
		log.Fatal("only docker-compose format v2 is supported")
	}

	//validate user
	if len(User) < 1 {
		log.Fatal("--user is required for the up command")
	}

	//prompt for password
	fmt.Printf("Password: ")
	passwd, _ := gopass.GetPasswd()
	var pass = string(passwd)

	//authenticate and get token
	var token = GetToken(User, pass)
	if Verbose {
		log.Printf("token obtained")
	}

	//iterate shipments
	for shipmentName, shipment := range harborCompose.Shipments {
		fmt.Println("Stopping " + shipmentName)

		if Verbose {
			log.Println("processing  " + shipmentName + "/" + shipment.Env)
			log.Println(shipment.Containers)
		}

		//override replicas
		shipment.Replicas = 0

		//update shipment level configuration
		UpdateShipment(shipmentName, shipment, token)

		//trigger shipment
		Trigger(shipmentName, shipment.Env)

		//TODO: delete shipment

		fmt.Println("done")
	}
}