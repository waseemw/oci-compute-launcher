# OCI Compute Launcher

## Description

Due to the lack of availability of A1 (Ampere) in Oracle Cloud Infrastructure, it is hard to manually create an instance with those CPUs,
since they come and go really quickly.

This program will, every 5 minutes, try to create an instance on each availability region (10 seconds delay between regions to further avoid
rate limiting). Once the instance is created, the instance ID will be printed and the program will stop.

## Instructions

1. Download the executable from the latest [release](https://github.com/waseemw/oci-compute-launcher/releases). You can also clone this
   repository and run `go build`, or just run `go run` in the last step
3. Create a `.env` file as a copy of `.env.example`
4. Set the credentials variables with the info you get by adding an API key in user settings in the OCI web interface
5. Set the Resource IDs variables, by going to each resource on OCI and copying the OCID
6. Set the Compute Config variables as you prefer. The values in `.env.example` are set for the always-free tier limits, given that you
   don't have any other Ampere-based Compute instances
7. Run the executable in the same directory as your `.env` file

## License

[MIT License](LICENSE)
