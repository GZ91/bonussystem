package flags

import "flag"

func ReadFlags() (string, string, string) {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	connectionStringDB := flag.String("d", "", "connection string for postgresql")
	billingSystem := flag.String("r", "", "address of the billing system")

	flag.Parse()
	return *addressServer, *connectionStringDB, *billingSystem
}
