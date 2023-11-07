package registry

import (
	"errors"
	"time"

	"golang.org/x/sys/windows/registry"
)

// Path to the key that we are adding values to. "LOCAL_MACHINE" is included in the actual call to open the key
var path = `SOFTWARE\WOW6432Node\Tanium\Tanium Client\Sensor Data\Tags\`

// Function that adds the tag we want. It creates the Registry value
func AddString(t string) error {
	tag := t
	currentTime := time.Now()
	//When Tanium adds a tag the value is a date formated like below. This does the same.
	date := currentTime.Format("1/02/2006 3:4:5 PM")
	keyString := "Added: " + date

	//Open the key first
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.WRITE|registry.SET_VALUE)
	if err != nil {
		//		fmt.Println("Error opening key: ", err)
		return errors.New("Error opening key. Try running as Administrator.")
	}
	//With the key open we can now add string value
	err = k.SetStringValue(tag, keyString)
	if err != nil {
		//		fmt.Println("Error setting string value: ", err)
		//		fmt.Println("\nTry running as Administrator", err)
		return errors.New("Error setting string value, try running as Administrator")
	}
	//Close the key after opening it
	k.Close()
	return nil
}

// Function that prints out the current Tanium tags on the machine. (Just reads the values in the key)
func GetKeyStrings() ([]string, error) {
	var values []string
	//We are opening the key with READ permissions, which allows this part to still execute when not running as Administrator
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		//fmt.Println("Error opening key: ", err)
		return nil, errors.New("Error opening key")
	}
	//The -1 is needed so that it doesn't return EOF in the error if you pass a regular number
	values, err = k.ReadValueNames(-1)
	if err != nil {
		//fmt.Println("Error reading values: ", err)
		return nil, errors.New("Error reading key values")
	}
	return values, nil
}
