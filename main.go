package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/v54/common"
	"github.com/oracle/oci-go-sdk/v54/core"
	"github.com/oracle/oci-go-sdk/v54/example/helpers"
	"github.com/oracle/oci-go-sdk/v54/identity"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	helpers.FatalIfError(godotenv.Load())
	keyBytes, err := os.ReadFile(os.Getenv("OCI_PRIVATE_KEY_PATH"))
	helpers.FatalIfError(err)
	tenancy := os.Getenv("OCI_TENANCY")
	user := os.Getenv("OCI_USER")
	region := os.Getenv("OCI_REGION")
	fingerPrint := os.Getenv("OCI_FINGERPRINT")
	subnet := os.Getenv("OCI_SUBNET")
	image := os.Getenv("OCI_IMAGE")
	shape := os.Getenv("COMPUTE_SHAPE")
	name := os.Getenv("COMPUTE_NAME")
	sshKey := os.Getenv("SSH_PUBLIC_KEY")
	privateIp := os.Getenv("NETWORK_PRIVATE_IP")

	cpuInt, err := strconv.Atoi(os.Getenv("COMPUTE_CPU"))
	memInt, err := strconv.Atoi(os.Getenv("COMPUTE_MEM"))
	volInt, err := strconv.Atoi(os.Getenv("COMPUTE_VOL"))
	fmt.Printf("Will launch %s, %d CPU, %d MEM, %d VOL \n\n", shape, cpuInt, memInt, volInt)

	helpers.FatalIfError(err)
	helpers.FatalIfError(err)
	cpu := float32(cpuInt)
	mem := float32(memInt)
	vol := int64(volInt)
	configProvider := common.NewRawConfigurationProvider(tenancy, user, region, fingerPrint, string(keyBytes), nil)
	identityClient, err := identity.NewIdentityClientWithConfigurationProvider(configProvider)
	helpers.FatalIfError(err)
	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	helpers.FatalIfError(err)
	domains, err := identityClient.ListAvailabilityDomains(context.Background(), identity.ListAvailabilityDomainsRequest{CompartmentId: &tenancy})
	helpers.FatalIfError(err)

	for {
		for _, item := range domains.Items {
			fmt.Printf("Trying on Availability Domain: %v\n", *item.Name)
			res, err := computeClient.LaunchInstance(context.Background(), core.LaunchInstanceRequest{
				LaunchInstanceDetails: core.LaunchInstanceDetails{
					CompartmentId: &tenancy,
					Shape:         &shape,
					CreateVnicDetails: &core.CreateVnicDetails{
						SubnetId:  &subnet,
						PrivateIp: &privateIp,
					},
					DisplayName: &name,
					ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
						Ocpus:       &cpu,
						MemoryInGBs: &mem,
					},
					AvailabilityDomain: item.Name,
					SourceDetails:      core.InstanceSourceViaImageDetails{ImageId: &image, BootVolumeSizeInGBs: &vol},
					Metadata: map[string]string{
						"ssh_authorized_keys": sshKey,
					},
				},
			})
			if err != nil {
				if strings.Contains(err.Error(), "Out of host capacity") {
					fmt.Printf("Out of host capacity\n")
					time.Sleep(10 * time.Second)
					continue
				}
				log.Fatal(err)
			}
			fmt.Printf("Launched! - Instance ID: %v\n", *res.Id)
			return
		}
		fmt.Printf("Retrying in 5 minutes...\n\n")
		time.Sleep(5 * time.Minute)
	}
}
