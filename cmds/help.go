package cmds;

import "fmt"

func Help() {
	fmt.Println("Welcome to all the users who wanted to refresh their memory regarding wgsl.")
	fmt.Println("The supported functions are")
	fmt.Println("To Initialise a new wgsl environment in the working directory w")
	fmt.Println("init    initialises a wgsl working directory which takes the email of the doctor storing it in the local folder for providing the results with the machine learning prediction")
	fmt.Println("To get the Image that we wanna test on ")
	fmt.Println("get  command that fetches the image that we wanna test from the secure location such as the servers to securely download it and test against")
	fmt.Println("To train the model")
	fmt.Println("train this command need to be run befor testing the machine learning on the new image that is provided")
	fmt.Println("to test the image")
	fmt.Println("test  this subcommand predicts the image/images via machine learning")
	fmt.Println("to get the result")
	fmt.Println("result sends all the results from running and testing the model is send to the email that is provided at the init stage")
}
